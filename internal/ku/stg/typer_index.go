package stg

import (
	"fmt"

	"github.com/mebyus/ku/internal/ku/ast"
	"github.com/mebyus/ku/internal/ku/enums/scok"
	"github.com/mebyus/ku/internal/ku/enums/symk"
)

func (t *Typer) index() {
	t.indexFuns()
}

func (t *Typer) indexFuns() {
	for i, f := range t.box.funs {
		t.indexFun(i, f)
	}
}

func (t *Typer) indexFun(i int, f *ast.Fun) {
	name := f.Name
	pin := f.Pin

	symbol := t.unit.Scope.Get(name)
	if symbol != nil {
		t.report(pin, fmt.Sprintf("symbol with name \"%s\" was already declared inside this unit", name))
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
