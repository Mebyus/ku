package typer

import (
	"fmt"

	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/enums/smk"
	"github.com/mebyus/ku/goku/compiler/typer/stg"
)

func (t *Typer) convConstSymbol(s *stg.Symbol) diag.Error {
	const debug = true

	c := t.box.Const(s.Aux)

	exp, err := t.evalConstExp(&t.unit.Scope, c.Exp)
	if err != nil {
		return err
	}

	if debug {
		fmt.Printf("const %s = %s\n", c.Name.Str, exp)
	}

	return nil
}

func (t *Typer) evalConstExp(scope *stg.Scope, exp ast.Exp) (stg.Exp, diag.Error) {
	switch e := exp.(type) {
	case ast.Nil:
		return nil, &diag.SimpleMessageError{
			Pin:  e.Pin,
			Text: "nil used in constant expression",
		}
	case ast.Symbol:
		return t.evalConstSymbolExp(scope, e)
	case ast.Integer:
		return t.ctx.Types.MakeInteger(e.Pin, e.Val), nil
	case ast.String:
		return t.ctx.Types.MakeString(e.Pin, e.Val), nil
	case ast.True:
	case ast.False:
	case ast.Rune:
	case ast.Binary:
		return t.evalConstBinaryExp(scope, e)
	default:
		panic(fmt.Sprintf("unexpected \"%s\" (=%d) expression (%T)", e.Kind(), e.Kind(), e))
	}

	// TODO: remove this return
	return nil, nil
}

func (t *Typer) evalConstSymbolExp(scope *stg.Scope, sym ast.Symbol) (stg.Exp, diag.Error) {
	name := sym.Name
	pin := sym.Pin

	s := scope.Lookup(name)
	if s == nil {
		return nil, &diag.SimpleMessageError{
			Pin:  pin,
			Text: fmt.Sprintf("type name \"%s\" refers to undefined symbol", name),
		}
	}
	if s.Kind != smk.Const {
		return nil, &diag.SimpleMessageError{
			Pin:  pin,
			Text: fmt.Sprintf("name \"%s\" refers to %s, not a constant", name, s.Kind),
		}
	}

	return s.Def.(stg.StaticValue).Exp, nil
}

func (t *Typer) evalConstBinaryExp(scope *stg.Scope, exp ast.Binary) (stg.Exp, diag.Error) {
	a, err := t.evalConstExp(scope, exp.A)
	if err != nil {
		return nil, err
	}
	b, err := t.evalConstExp(scope, exp.B)
	if err != nil {
		return nil, err
	}
	if a.Type() != b.Type() {
		return nil, &diag.SimpleMessageError{
			Pin:  a.Span().Pin,
			Text: fmt.Sprintf("incompatible operand types (%s and %s) in binary expression", a.Type(), b.Type()),
		}
	}
	return nil, nil
}
