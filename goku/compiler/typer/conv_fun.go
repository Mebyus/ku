package typer

import (
	"fmt"

	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/enums/sck"
	"github.com/mebyus/ku/goku/compiler/enums/smk"
	"github.com/mebyus/ku/goku/compiler/typer/stg"
)

func (t *Typer) convFunSymbol(s *stg.Symbol) diag.Error {
	fun := t.box.Fun(s.Aux)
	def := &stg.Fun{}
	s.Def = def
	scope := &def.Body.Scope
	scope.Init(sck.Node, &t.unit.Scope)

	if def.Never && fun.Signature.Result != nil {
		panic("invalid signature")
	}
	def.Never = fun.Signature.Never

	if fun.Signature.Result != nil {
		result, err := t.ctx.Types.Lookup(&t.unit.Scope, fun.Signature.Result)
		if err != nil {
			return err
		}
		def.Result = result
	}

	var params []*stg.Type
	if len(fun.Signature.Params) != 0 {
		params = make([]*stg.Type, 0, len(fun.Signature.Params))
	}
	for _, p := range fun.Signature.Params {
		name := p.Name.Str
		pin := p.Name.Pin

		if scope.Has(name) {
			return &diag.SimpleMessageError{
				Pin:  pin,
				Text: fmt.Sprintf("multiple parameters with name \"%s\" in function \"%s\"", name, s.Name),
			}
		}
		ps := scope.Alloc(smk.Param, name, pin)

		typ, err := t.ctx.Types.Lookup(&t.unit.Scope, p.Type)
		if err != nil {
			return err
		}
		ps.Type = typ

		params = append(params, typ)
	}
	def.Params = params

	return nil
}

func (t *Typer) convMethodSymbol(s *stg.Symbol) diag.Error {
	return nil
}
