package stg

import (
	"fmt"

	"github.com/mebyus/ku/internal/ku/ast"
	"github.com/mebyus/ku/internal/ku/enums/scok"
	"github.com/mebyus/ku/internal/ku/enums/symk"
	"github.com/mebyus/ku/internal/ku/enums/typk"
	"github.com/mebyus/ku/internal/ku/sx"
)

func (t *Typer) index() {
	t.indexImports()
	t.indexConsts()
	t.indexStubs()
	t.indexFuns()
	t.indexTypes()
}

func (t *Typer) indexImports() {
	for i, p := range t.box.imps {
		t.indexImport(i, &p)
	}
}

func (t *Typer) indexConsts() {
	for i, c := range t.box.consts {
		t.indexConst(i, &c)
	}
}

func (t *Typer) indexFuns() {
	for i, f := range t.box.funs {
		t.indexFun(i, &f)
	}
}

func (t *Typer) indexStubs() {
	for i, s := range t.box.stubs {
		t.indexStub(i, &s)
	}
}

func (t *Typer) indexTypes() {
	for i, s := range t.box.types {
		t.indexType(i, &s)
	}
}

func (t *Typer) indexImport(i int, p *ast.Import) {
	name := p.Name
	pin := p.Pin

	unit, ok := t.common.Units[p.Path]
	if !ok || unit == nil {
		panic("imported unit not found")
	}

	symbol := t.unit.Scope.Get(name)
	if symbol != nil {
		t.report(pin, sx.ETypDupSymbol, fmt.Sprintf("symbol with name \"%s\" was already declared inside this unit", name))
		return
	}

	symbol = t.unit.Scope.New(symk.Import, name, pin)
	symbol.Def = &Import{Unit: unit}
}

func (t *Typer) indexConst(i int, c *ast.Const) {
	name := c.Name
	pin := c.Pin

	symbol := t.unit.Scope.Get(name)
	if symbol != nil {
		t.report(pin, sx.ETypDupSymbol, fmt.Sprintf("symbol with name \"%s\" was already declared inside this unit", name))
		return
	}

	symbol = t.unit.Scope.New(symk.Const, name, pin)
	symbol.Aux = uint32(i)

	t.unit.Consts = append(t.unit.Consts, symbol)
}

func (t *Typer) indexStub(i int, s *ast.FunStub) {
	name := s.Name
	pin := s.Pin

	symbol := t.unit.Scope.Get(name)
	if symbol != nil {
		t.report(pin, sx.ETypDupSymbol, fmt.Sprintf("symbol with name \"%s\" was already declared inside this unit", name))
		return
	}

	symbol = t.unit.Scope.New(symk.Fun, name, pin)
	symbol.Aux = uint32(i)
	symbol.Def = t.newFunDef()
	symbol.Flags |= SymbolStub

	t.unit.Funs = append(t.unit.Funs, symbol)
}

func (t *Typer) indexFun(i int, f *ast.Fun) {
	name := f.Name
	pin := f.Pin

	symbol := t.unit.Scope.Get(name)
	if symbol != nil {
		t.report(pin, sx.ETypDupSymbol, fmt.Sprintf("symbol with name \"%s\" was already declared inside this unit", name))
		return
	}

	symbol = t.unit.Scope.New(symk.Fun, name, pin)
	symbol.Aux = uint32(i)
	symbol.Def = t.newFunDef()

	// if f.Pub {
	// 	symbol.Flags |= SymbolPublic
	// }
	// if f.Export {
	// 	symbol.Flags |= SymbolExport
	// 	t.unit.Export = append(t.unit.Export, symbol)
	// }

	t.unit.Funs = append(t.unit.Funs, symbol)
}

func (t *Typer) newFunDef() *FunDef {
	def := &FunDef{}
	def.Body.Scope.Init(scok.Node, &t.unit.Scope)
	return def
}

func (t *Typer) indexType(i int, s *ast.Type) {
	name := s.Name
	pin := s.Pin

	symbol := t.unit.Scope.Get(name)
	if symbol != nil {
		t.report(pin, sx.ETypDupSymbol, fmt.Sprintf("symbol with name \"%s\" was already declared inside this unit", name))
		return
	}

	symbol = t.unit.Scope.New(symk.Type, name, pin)
	symbol.Aux = uint32(i)
	symbol.Def = t.newTypeDef(symbol)
}

func (t *Typer) newTypeDef(s *Symbol) *Type {
	def := &Custom{Symbol: s}
	typ := &Type{
		Def:  def,
		Kind: typk.Custom,
	}
	return typ
}
