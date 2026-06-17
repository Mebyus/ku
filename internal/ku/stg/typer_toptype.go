package stg

import (
	"fmt"

	"github.com/mebyus/ku/internal/ku/ast"
	"github.com/mebyus/ku/internal/ku/enums/symk"
	"github.com/mebyus/ku/internal/ku/enums/typk"
)

// assigns types for top-level variables and function params and results
func (t *Typer) toptype() {
	for _, s := range t.unit.Funs {
		if s.IsStub() {
			t.typeFun(s.Def.(*FunDef), &t.box.stubs[s.Aux].Sig)
		} else {
			t.typeFun(s.Def.(*FunDef), &t.box.funs[s.Aux].Sig)
		}
	}
}

func (t *Typer) typeFun(def *FunDef, sig *ast.Signature) {
	if def.Never && sig.Result != nil {
		panic("invalid signature")
	}
	def.Never = sig.Never

	if sig.Result != nil {
		result := t.LookupType(&t.unit.Scope, sig.Result)
		if result.Kind != typk.Void {
			def.Result = result
		}
	}

	scope := &def.Body.Scope
	var inputs []*Type
	var params []*Symbol
	if len(sig.Params) != 0 {
		inputs = make([]*Type, 0, len(sig.Params))
		params = make([]*Symbol, 0, len(sig.Params))
	}
	for _, p := range sig.Params {
		name := p.Name
		pin := p.Pin

		symbol := scope.Get(name)
		if symbol != nil {
			t.report(pin, fmt.Sprintf("function already has parameter named \"%s\"", name))
			continue
		}

		typ := t.LookupType(&t.unit.Scope, p.Type)
		symbol = scope.New(symk.Param, name, pin)
		symbol.Type = typ

		inputs = append(inputs, typ)
		params = append(params, symbol)
	}
	def.Inputs = inputs
	def.Params = params
}
