package typer

import (
	"fmt"

	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/enums/smk"
	"github.com/mebyus/ku/goku/compiler/typer/stg"
)

func (t *Typer) convConstSymbol(s *stg.Symbol) diag.Error {
	c := t.box.Const(s.Aux)

	_, err := t.evalConstExp(&t.unit.Scope, c.Exp)
	if err != nil {
		return err
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
	case ast.String:
	case ast.True:
	case ast.False:
	case ast.Rune:
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

	return s.Def.(stg.Exp), nil
}
