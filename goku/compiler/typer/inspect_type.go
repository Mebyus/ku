package typer

import (
	"fmt"

	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/enums/smk"
	"github.com/mebyus/ku/goku/compiler/typer/stg"
)

func (t *Typer) inspectTypeSymbol(s *stg.Symbol) diag.Error {
	spec := t.box.Type(s.Aux).Spec

	err := t.inspectCustomTypeSpec(spec)
	if err != nil {
		return err
	}

	kind, ok := t.ins.links[s]
	if ok && kind == LinkDirect {
		return &diag.SimpleMessageError{
			Pin:  s.Pin,
			Text: fmt.Sprintf("type \"%s\" definition directly references itself", s.Name),
		}
	}

	return nil
}

func (t *Typer) inspectCustomTypeSpec(spec ast.TypeSpec) diag.Error {
	switch p := spec.(type) {
	case ast.Void:
		return nil
	case ast.Struct:
		return t.inspectFields(p.Fields)
	case ast.Bag:
		fmt.Printf("WARN: bag type specifier not implemented\n")
		return nil
	case ast.Enum:
		return t.inspectEnum(p)
	case ast.Union:
		return t.inspectFields(p.Fields)
	}

	return t.inspectType(spec)
}

func (t *Typer) inspectFields(fields []ast.Field) diag.Error {
	for _, f := range fields {
		err := t.inspectType(f.Type)
		if err != nil {
			return err
		}
	}
	return nil
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

	unit := m.Def.(stg.SymDefUnit).Unit
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

func (t *Typer) inspectTuple(p ast.Tuple) diag.Error {
	for _, p := range p.Types {
		err := t.inspectType(p)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Typer) inspectEnum(p ast.Enum) diag.Error {
	for _, e := range p.Entries {
		if e.Exp == nil {
			continue
		}

		err := t.inspectConstExp(e.Exp)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Typer) inspectType(spec ast.TypeSpec) diag.Error {
	switch p := spec.(type) {
	case ast.VoidPointer, ast.VoidRef:
		return nil
	case ast.TypeName:
		return t.linkTypeName(p)
	case ast.TypeFullName:
		return t.inspectTypeFullName(p)
	case ast.Pointer:
		return t.linkPointer(p)
	case ast.Ref:
		return t.linkRef(p)
	case ast.ArrayPointer:
		return t.linkArrayPointer(p)
	case ast.Span:
		return t.linkSpan(p)
	case ast.CapBuf:
		return t.linkCapBuf(p)
	case ast.Array:
		return t.linkArray(p)
	case ast.ArrayRef:
		return t.linkArrayRef(p)
	default:
		panic(fmt.Sprintf("unexpected \"%s\" (=%d) type specifier (%T)", p.Kind(), p.Kind(), p))
	}
}

func (t *Typer) linkPointer(p ast.Pointer) diag.Error {
	k := t.ins.indirect()
	err := t.inspectType(p.Type)
	t.ins.restore(k)
	return err
}

func (t *Typer) linkArrayRef(p ast.ArrayRef) diag.Error {
	k := t.ins.indirect()
	err := t.inspectType(p.Type)
	t.ins.restore(k)
	return err
}

func (t *Typer) linkArrayPointer(p ast.ArrayPointer) diag.Error {
	k := t.ins.indirect()
	err := t.inspectType(p.Type)
	t.ins.restore(k)
	return err
}

func (t *Typer) linkRef(p ast.Ref) diag.Error {
	k := t.ins.indirect()
	err := t.inspectType(p.Type)
	t.ins.restore(k)
	return err
}

func (t *Typer) linkSpan(c ast.Span) diag.Error {
	k := t.ins.indirect()
	err := t.inspectType(c.Type)
	t.ins.restore(k)
	return err
}

func (t *Typer) linkCapBuf(c ast.CapBuf) diag.Error {
	k := t.ins.indirect()
	err := t.inspectType(c.Type)
	t.ins.restore(k)
	return err
}

func (t *Typer) linkArray(a ast.Array) diag.Error {
	if a.Size != nil {
		err := t.inspectConstExp(a.Size)
		if err != nil {
			return err
		}
	}

	return t.inspectType(a.Type)
}
