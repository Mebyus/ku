package typer

import (
	"fmt"

	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/enums/sck"
	"github.com/mebyus/ku/goku/compiler/enums/smk"
	"github.com/mebyus/ku/goku/compiler/typer/stg"
)

func (t *Typer) convFunSymbol(s *stg.Symbol) diag.Error {
	fun := t.box.Fun(s.Aux)
	def, err := t.createFun(s.Name, fun.Signature)
	if err != nil {
		return err
	}

	s.Def = def
	return nil
}

func (t *Typer) convMethodSymbol(s *stg.Symbol) diag.Error {
	m := t.box.Method(s.Aux)
	def, err := t.createFun(s.Name, m.Signature)
	if err != nil {
		return err
	}

	r := m.Receiver
	var spec ast.TypeSpec
	switch r.Kind {
	case ast.ReceiverVal:
		spec = ast.TypeName{Name: r.Name}
	case ast.ReceiverPtr:
		spec = ast.Pointer{Type: ast.TypeName{Name: r.Name}}
	case ast.ReceiverRef:
		spec = ast.Ref{Type: ast.TypeName{Name: r.Name}}
	default:
		panic(fmt.Sprintf("unexpected receiver kind (=%d)", r.Kind))
	}

	typ, err := t.ctx.Types.Lookup(&t.unit.Scope, spec)
	if err != nil {
		return err
	}
	def.Receiver = typ


	const rname = "g"
	if def.Body.Scope.Has(rname) {
		return &diag.SimpleMessageError{
			Pin:  r.Name.Pin,
			Text: fmt.Sprintf("multiple parameters with name \"%s\" in function \"%s\"", rname, s.Name),
		}
	}
	rs := def.Body.Scope.Alloc(smk.Receiver, rname, r.Name.Pin)
	rs.Type = typ

	s.Def = def
	return nil
}

func (t *Typer) createFun(sname string, sig ast.Signature) (*stg.Fun, diag.Error) {
	def := &stg.Fun{}
	scope := &def.Body.Scope
	scope.Init(sck.Node, &t.unit.Scope)

	if def.Never && sig.Result != nil {
		panic("invalid signature")
	}
	def.Never = sig.Never

	if sig.Result != nil {
		result, err := t.ctx.Types.Lookup(&t.unit.Scope, sig.Result)
		if err != nil {
			return nil, err
		}
		def.Result = result
	}

	var params []*stg.Type
	if len(sig.Params) != 0 {
		params = make([]*stg.Type, 0, len(sig.Params))
	}
	for _, p := range sig.Params {
		name := p.Name.Str
		pin := p.Name.Pin

		if scope.Has(name) {
			return nil, &diag.SimpleMessageError{
				Pin:  pin,
				Text: fmt.Sprintf("multiple parameters with name \"%s\" in function \"%s\"", name, sname),
			}
		}
		ps := scope.Alloc(smk.Param, name, pin)

		typ, err := t.ctx.Types.Lookup(&t.unit.Scope, p.Type)
		if err != nil {
			return nil, err
		}
		ps.Type = typ

		params = append(params, typ)
	}
	def.Params = params

	return def, nil
}
