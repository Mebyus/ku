package typer

import (
	"fmt"

	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/enums/smk"
)

// Expression should evaluate to constant, report error otherwise.
//
// Report error upon encountering variable symbol.
func (t *Typer) inspectConstExp(exp ast.Exp) diag.Error {
	switch e := exp.(type) {
	case ast.Integer, ast.String, ast.Rune, ast.True, ast.False, ast.Float:
		return nil
	case ast.Symbol:
		return t.linkConstSymbol(e)
	case ast.Unary:
		return t.inspectConstUnaryExp(e)
	case ast.Binary:
		return t.inspectConstBinaryExp(e)
	case ast.Cast:
		return t.inspectConstCastExp(e)
	case ast.List:
		return t.inspectConstListExp(e)
	case ast.Paren:
		return t.inspectConstExp(e.Exp)
	default:
		panic(fmt.Sprintf("unexpected \"%s\" (=%d) expression (%T)", e.Kind(), e.Kind(), e))
	}
}

func (t *Typer) inspectConstListExp(l ast.List) diag.Error {
	for _, e := range l.Exps {
		err := t.inspectConstExp(e)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Typer) inspectConstCastExp(c ast.Cast) diag.Error {
	err := t.inspectVarType(c.Type)
	if err != nil {
		return err
	}
	return t.inspectConstExp(c.Exp)
}

func (t *Typer) inspectConstUnaryExp(exp ast.Unary) diag.Error {
	return t.inspectConstExp(exp.Exp)
}

func (t *Typer) inspectConstBinaryExp(exp ast.Binary) diag.Error {
	err := t.inspectConstExp(exp.A)
	if err != nil {
		return err
	}
	return t.inspectConstExp(exp.B)
}

func (t *Typer) linkConstSymbol(sym ast.Symbol) diag.Error {
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
			Text: fmt.Sprintf("name \"%s\" refers to %s, not a constant", name, s.Kind),
		}
	}

	t.ins.link(s)
	return nil
}
