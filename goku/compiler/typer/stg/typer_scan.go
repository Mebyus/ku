package stg

import (
	"fmt"

	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/enums/sck"
	"github.com/mebyus/ku/goku/compiler/enums/smk"
)

// DepSet is a hepler object for collecting dependencies of unit-level symbols
// on other unit-level symbols.
type DepSet struct {
	// Contains m of symbols referenced by the symbol being scanned.
	//
	// Currently we store symbol raw ids to avoid GC scans of constructed
	// dependency graph.
	m map[ /* symbol id */ uint]struct{}
}

func (g *DepSet) init() {
	g.m = make(map[uint]struct{})
}

func (g *DepSet) reset() {
	clear(g.m)
}

func (g *DepSet) has(s *Symbol) bool {
	_, ok := g.m[s.RawID()]
	return ok
}

// Take returns list of all dependencies found during scan and resets
// internal state afterward. After calling this method DepSet can be
// used again for scanning another symbol.
func (g *DepSet) take() []uint {
	if len(g.m) == 0 {
		return nil
	}

	list := make([]uint, 0, len(g.m))
	for s := range g.m {
		list = append(list, s)
	}

	g.reset()
	return list
}

// Add dependency between the given symbol and the one currently under inspection.
func (g *DepSet) add(s *Symbol) {
	if s.Scope.Kind != sck.Unit {
		// do not link global symbols, they are defined implicitly
		// before everything else
		return
	}

	g.m[s.RawID()] = struct{}{}
}

func (t *Typer) scanConstSymbols() {
	for i, s := range t.consts {
		err := t.scanConstSymbol(i, s)
		if err != nil {
			t.report(err)
		}

		_ = t.deps.take()
	}
}

func (t *Typer) scanVarSymbols() {
	for i, s := range t.vars {
		err := t.scanVarSymbol(i, s)
		if err != nil {
			t.report(err)
		}

		// TODO: we can avoid map insteraction when scanning
		// variables, just need to verify type spec and init expression
		_ = t.deps.take()
	}
}

func (t *Typer) scanTypeSymbols() {
	for i, s := range t.types {
		err := t.scanTypeSymbol(i, s)
		if err != nil {
			t.report(err)
		}

		_ = t.deps.take()
	}
}

func (t *Typer) scanConstSymbol(i int, s *Symbol) diag.Error {
	c := t.box.consts[i]

	if c.Type != nil {
		err := t.scanType(c.Type)
		if err != nil {
			return err
		}
	}

	err := t.scanExp(c.Exp)
	if err != nil {
		return err
	}

	if t.deps.has(s) {
		return &diag.SimpleMessageError{
			Pin:  s.Pin,
			Text: fmt.Sprintf("constant \"%s\" definition references itself", s.Name),
		}
	}

	return nil
}

func (t *Typer) scanVarSymbol(i int, s *Symbol) diag.Error {
	v := t.box.vars[i]
	_ = v

	err := t.scanType(v.Type)
	if err != nil {
		return err
	}

	err = t.scanExp(v.Exp)
	if err != nil {
		return err
	}

	return nil
}

func (t *Typer) scanTypeSymbol(i int, s *Symbol) diag.Error {
	typ := t.box.types[i]
	_ = typ

	err := t.scanCustomTypeSpec(typ.Spec)
	if err != nil {
		return err
	}

	return nil
}

// Expression should evaluate to constant, report error otherwise.
//
// Report error upon encountering variable symbol.
func (t *Typer) scanExp(exp ast.Exp) diag.Error {
	switch e := exp.(type) {
	case ast.Integer, ast.String, ast.Rune, ast.True, ast.False, ast.Float:
		return nil
	case ast.Symbol:
		return t.scanSymbolExp(e)
	case ast.Unary:
		return t.scanExp(e.Exp)
	case ast.Binary:
		return t.scanBinary(e)
	case ast.Cast:
		return t.scanCast(e)
	case ast.List:
		return t.scanList(e)
	case ast.Paren:
		return t.scanExp(e.Exp)
	default:
		panic(fmt.Sprintf("unexpected \"%s\" (=%d) expression (%T)", e.Kind(), e.Kind(), e))
	}
}

func (t *Typer) scanList(l ast.List) diag.Error {
	for _, e := range l.Exps {
		err := t.scanExp(e)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Typer) scanCast(c ast.Cast) diag.Error {
	err := t.scanType(c.Type)
	if err != nil {
		return err
	}
	return t.scanExp(c.Exp)
}

func (t *Typer) scanBinary(exp ast.Binary) diag.Error {
	err := t.scanExp(exp.A)
	if err != nil {
		return err
	}
	return t.scanExp(exp.B)
}

func (t *Typer) scanSymbolExp(sym ast.Symbol) diag.Error {
	name := sym.Name
	s := t.unit.Scope.Lookup(name)
	if s == nil {
		return &diag.SimpleMessageError{
			Pin:  sym.Pin,
			Text: fmt.Sprintf("name \"%s\" refers to undefined symbol", name),
		}
	}
	if s.Kind != smk.Const {
		return &diag.SimpleMessageError{
			Pin:  sym.Pin,
			Text: fmt.Sprintf("name \"%s\" refers to %s which is not allowed in compile-time expression", name, s.Kind),
		}
	}

	t.deps.add(s)
	return nil
}

func (t *Typer) scanCustomTypeSpec(spec ast.TypeSpec) diag.Error {
	switch p := spec.(type) {
	case ast.Void:
		return nil
	case ast.Struct:
		return t.scanFields(p.Fields)
	case ast.Bag:
		fmt.Printf("WARN: bag type specifier not implemented\n")
		return nil
	case ast.Enum:
		return t.scanEnum(p)
	case ast.Union:
		return t.scanFields(p.Fields)
	}

	return t.scanType(spec)
}

func (t *Typer) scanFields(fields []ast.Field) diag.Error {
	for _, f := range fields {
		err := t.scanType(f.Type)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Typer) scanTypeFullName(p ast.TypeFullName) diag.Error {
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

	unit := m.Def.(*SymDefUnit).Unit
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

func (t *Typer) scanTypeName(p ast.TypeName) diag.Error {
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

	t.deps.add(s)
	return nil
}

func (t *Typer) scanEnum(p ast.Enum) diag.Error {
	for _, e := range p.Entries {
		if e.Exp == nil {
			continue
		}

		err := t.scanExp(e.Exp)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Typer) scanType(spec ast.TypeSpec) diag.Error {
	switch p := spec.(type) {
	case ast.VoidPointer, ast.VoidRef:
		return nil
	case ast.TypeName:
		return t.scanTypeName(p)
	case ast.TypeFullName:
		return t.scanTypeFullName(p)
	case ast.Pointer:
		return t.scanType(p.Type)
	case ast.Ref:
		return t.scanType(p)
	case ast.ArrayPointer:
		return t.scanType(p.Type)
	case ast.Span:
		return t.scanType(p.Type)
	case ast.CapBuf:
		return t.scanType(p.Type)
	case ast.Array:
		return t.scanArray(p)
	case ast.ArrayRef:
		return t.scanType(p.Type)
	case ast.Map:
		return t.scanMap(p)
	default:
		panic(fmt.Sprintf("unexpected \"%s\" (=%d) type specifier (%T)", p.Kind(), p.Kind(), p))
	}
}

func (t *Typer) scanMap(m ast.Map) diag.Error {
	err := t.scanType(m.Key)
	if err != nil {
		return err
	}
	return t.scanType(m.Value)
}

func (t *Typer) scanArray(a ast.Array) diag.Error {
	if a.Size != nil {
		err := t.scanExp(a.Size)
		if err != nil {
			return err
		}
	}

	return t.scanType(a.Type)
}
