package typer

import (
	"fmt"

	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/enums/smk"
	"github.com/mebyus/ku/goku/compiler/enums/tpk"
	"github.com/mebyus/ku/goku/compiler/typer/stg"
)

func (t *Typer) inspectTypeSymbol(s *stg.Symbol) diag.Error {
	spec := t.box.Type(s.Aux).Spec

	custom := &stg.Custom{
		Symbol:  s,
		Methods: t.methodsByReceiver[s],
	}
	def := &stg.Type{
		Def:  custom,
		Kind: tpk.Custom,
	}
	s.Def = def

	// TODO: probably need to add new type to index
	// t.ctx.Types.

	var err diag.Error
	switch p := spec.(type) {
	case ast.TypeName:
		return t.linkTypeName(p)
	case ast.Chunk:
		return t.linkChunk(p)
	case ast.Pointer:
		return t.linkPointer(p)
	case ast.Array:
		return t.linkArray(p)
	case ast.AnyPointer:
		return nil
	case ast.Struct:
		// TODO: check for unique names among fields + methods
		return t.inspectStructFields(p.Fields)
	case ast.Bag:
		fmt.Printf("WARN: bag type specifier not implemented (%s %d %T)\n", s.Name, s.Aux, spec)
	case ast.Union:
		fmt.Printf("WARN: union type specifier not implemented (%s %d %T)\n", s.Name, s.Aux, spec)
	default:
		panic(fmt.Sprintf("unexpected \"%s\" (=%d) type specifier (%T)", p.Kind(), p.Kind(), p))
	}
	return err
}

func (t *Typer) inspectStructFields(fields []ast.Field) diag.Error {
	for _, f := range fields {
		err := t.inspectStructField(f)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Typer) inspectStructField(field ast.Field) diag.Error {
	switch p := field.Type.(type) {
	case ast.TypeName:
		return t.linkTypeName(p)
	case ast.TypeFullName:
		return t.inspectTypeFullName(p)
	case ast.Pointer:
		return t.linkPointer(p)
	case ast.ArrayPointer:
		return t.linkArrayPointer(p)
	case ast.AnyPointer:
		return nil
	case ast.Array:
		return t.linkArray(p)
	case ast.Chunk:
		return t.linkChunk(p)
	default:
		panic(fmt.Sprintf("unexpected \"%s\" (=%d) type specifier (%T)", p.Kind(), p.Kind(), p))
	}
}

func (t *Typer) inspectTypeFullName(p ast.TypeFullName) diag.Error {
	// just check import + used type name
	// we do not need to link them, because they already come
	// from units with lower rank

	iname := p.Import.Str
	m := t.unit.Scope.Lookup(iname)
	if m == nil {
		return &diag.SimpleMessageError{
			Pin:  p.Name.Pin,
			Text: fmt.Sprintf("name \"%s\" refers to undefined symbol", iname),
		}
	}
	if m.Kind != smk.Import {
		return &diag.SimpleMessageError{
			Pin:  p.Name.Pin,
			Text: fmt.Sprintf("name \"%s\" refers to %s, not an import", iname, m.Kind),
		}
	}

	unit := m.Def.(*stg.Unit)
	name := p.Name.Str
	s := unit.Scope.Lookup(name)
	if s == nil {
		return &diag.SimpleMessageError{
			Pin:  p.Name.Pin,
			Text: fmt.Sprintf("name \"%s\" refers to undefined symbol", name),
		}
	}
	if s.Kind != smk.Type {
		return &diag.SimpleMessageError{
			Pin:  p.Name.Pin,
			Text: fmt.Sprintf("name \"%s\" refers to %s, not a type", name, s.Kind),
		}
	}
	if !s.IsPublic() {
		return &diag.SimpleMessageError{
			Pin:  p.Name.Pin,
			Text: fmt.Sprintf("type \"%s\" is not public", name),
		}
	}

	return nil
}

func (t *Typer) linkTypeName(p ast.TypeName) diag.Error {
	name := p.Name.Str
	s := t.unit.Scope.Lookup(name)
	if s == nil {
		return &diag.SimpleMessageError{
			Pin:  p.Name.Pin,
			Text: fmt.Sprintf("type name \"%s\" refers to undefined symbol", name),
		}
	}
	if s.Kind != smk.Type {
		return &diag.SimpleMessageError{
			Pin:  p.Name.Pin,
			Text: fmt.Sprintf("name \"%s\" refers to %s, not a type", name, s.Kind),
		}
	}

	t.ins.link(s)
	return nil
}

func (t *Typer) linkPointer(p ast.Pointer) diag.Error {
	k := t.ins.indirect()

	var err diag.Error
	switch p := p.Type.(type) {
	case ast.TypeName:
		err = t.linkTypeName(p)
	case ast.Pointer:
		err = t.linkPointer(p)
	case ast.AnyPointer:
		// do nothing
	default:
		panic(fmt.Sprintf("unexpected \"%s\" (=%d) type specifier (%T)", p.Kind(), p.Kind(), p))
	}

	t.ins.restore(k)
	return err
}

func (t *Typer) linkArrayPointer(p ast.ArrayPointer) diag.Error {
	// TODO: do we need to unify this function with linkPointer?
	// they should behave identically anyway
	k := t.ins.indirect()

	var err diag.Error
	switch p := p.Type.(type) {
	case ast.TypeName:
		err = t.linkTypeName(p)
	case ast.Pointer:
		err = t.linkPointer(p)
	case ast.AnyPointer:
		// do nothing
	default:
		panic(fmt.Sprintf("unexpected \"%s\" (=%d) type specifier (%T)", p.Kind(), p.Kind(), p))
	}

	t.ins.restore(k)
	return err
}


func (t *Typer) linkRef(p ast.Ref) diag.Error {
	// TODO: do we need to unify this function with linkPointer?
	// they should behave identically anyway
	k := t.ins.indirect()

	var err diag.Error
	switch p := p.Type.(type) {
	case ast.TypeName:
		err = t.linkTypeName(p)
	case ast.Pointer:
		err = t.linkPointer(p)
	case ast.AnyPointer:
		// do nothing
	default:
		panic(fmt.Sprintf("unexpected \"%s\" (=%d) type specifier (%T)", p.Kind(), p.Kind(), p))
	}

	t.ins.restore(k)
	return err
}

func (t *Typer) linkChunk(c ast.Chunk) diag.Error {
	k := t.ins.indirect()

	var err diag.Error
	switch p := c.Type.(type) {
	case ast.TypeName:
		err = t.linkTypeName(p)
	case ast.Pointer:
		err = t.linkPointer(p)
	case ast.AnyPointer:
		// do nothing
	default:
		panic(fmt.Sprintf("unexpected \"%s\" (=%d) type specifier (%T)", p.Kind(), p.Kind(), p))
	}

	t.ins.restore(k)
	return err
}

func (t *Typer) linkArray(a ast.Array) diag.Error {
	err := t.inspectConstExp(a.Size)
	if err != nil {
		return err
	}

	switch p := a.Type.(type) {
	case ast.TypeName:
		return t.linkTypeName(p)
	case ast.Pointer:
		return t.linkPointer(p)
	case ast.AnyPointer:
		return nil
	default:
		panic(fmt.Sprintf("unexpected \"%s\" (=%d) type specifier (%T)", p.Kind(), p.Kind(), p))
	}
}
