package typer

import (
	"fmt"

	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/enums/bok"
	"github.com/mebyus/ku/goku/compiler/enums/smk"
	"github.com/mebyus/ku/goku/compiler/enums/tpk"
	"github.com/mebyus/ku/goku/compiler/typer/stg"
)

func (t *Typer) convConstSymbol(s *stg.Symbol) diag.Error {
	const debug = true

	c := t.box.Const(s.Aux)
	name := c.Name.Str

	exp, err := t.evalConstExp(&t.unit.Scope, c.Exp)
	if err != nil {
		return err
	}

	var typ *stg.Type
	if c.Type == nil {
		typ = exp.Type()
	} else {
		var err diag.Error
		typ, err = t.ctx.Types.Lookup(&t.unit.Scope, c.Type)
		if err != nil {
			return err
		}
		err = checkConstValueType(typ, exp)
		if err != nil {
			return err
		}
		// TODO: convert exp to the same type as declared
	}

	s.Def = stg.StaticValue{Exp: exp}

	if debug {
		fmt.Printf("const %s: %s = %s\n", name, typ, exp)
	}

	return nil
}

func checkConstValueType(want *stg.Type, exp stg.Exp) diag.Error {
	t := exp.Type()

	switch want.Kind {
	case tpk.Custom:
		c := want.Def.(*stg.Custom)
		if c.Type.Kind == tpk.Integer && t.Kind == tpk.Integer && t.Size == 0 {
			return nil
		}
	default:
		panic(fmt.Sprintf("unexpected \"%s\" (=%d) type (%T)", want.Kind, want.Kind, want.Def))
	}

	return &diag.SimpleMessageError{
		Pin:  exp.Span().Pin,
		Text: fmt.Sprintf("declared constant type %s and definition value type %s are incompatible", want, exp.Type()),
	}
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
		// TODO: make separate type kind for runes
		return t.ctx.Types.MakeInteger(e.Pin, e.Val), nil
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

	typ := a.Type()
	switch typ.Kind {
	case tpk.Integer:
		if typ.Size != 0 {
			panic(fmt.Sprintf("static sized (size=%d) integers not implemented", typ.Size))
		}

		switch exp.Op.Kind {
		case bok.Add:
			return t.addIntegers(a.(stg.Integer), b.(stg.Integer)), nil
		default:
			panic(fmt.Sprintf("unexpected \"%s\" (=%d) binary operator", exp.Op.Kind, exp.Op.Kind))
		}
	default:
		panic(fmt.Sprintf("unexpected \"%s\" (=%d) type (%T)", typ.Kind, typ.Kind, typ.Def))
	}
}

func (t *Typer) addIntegers(a, b stg.Integer) stg.Integer {
	switch {
	case a.Neg && b.Neg:
		return t.ctx.Types.MakeNegInteger(a.Pin, a.Val+b.Val)
	case a.Neg && !b.Neg:
		if b.Val > a.Val {
			return t.ctx.Types.MakeNegInteger(a.Pin, b.Val-a.Val)
		}
		return t.ctx.Types.MakeInteger(a.Pin, a.Val-b.Val)
	case !a.Neg && b.Neg:
		if a.Val > b.Val {
			return t.ctx.Types.MakeNegInteger(a.Pin, a.Val-b.Val)
		}
		return t.ctx.Types.MakeInteger(a.Pin, b.Val-a.Val)
	case !a.Neg && !b.Neg:
		return t.ctx.Types.MakeInteger(a.Pin, a.Val+b.Val)
	default:
		panic("impossible condition")
	}
}
