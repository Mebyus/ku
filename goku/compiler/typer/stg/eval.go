package stg

import (
	"fmt"

	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/enums/bok"
	"github.com/mebyus/ku/goku/compiler/enums/smk"
	"github.com/mebyus/ku/goku/compiler/enums/tpk"
)

func (s *Scope) EvalConstExp(exp ast.Exp) (Exp, diag.Error) {
	switch e := exp.(type) {
	case ast.Nil:
		return nil, &diag.SimpleMessageError{
			Pin:  e.Pin,
			Text: "nil used in constant expression",
		}
	case ast.Symbol:
		return s.evalConstSymbolExp(e)
	case ast.Integer:
		return s.Types.MakeInteger(e.Pin, e.Val), nil
	case ast.String:
		return s.Types.MakeString(e.Pin, e.Val), nil
	case ast.True:
	case ast.False:
	case ast.Rune:
		// TODO: make separate type kind for runes
		return s.Types.MakeInteger(e.Pin, e.Val), nil
	case ast.Binary:
		return s.evalConstBinaryExp(e)
	default:
		panic(fmt.Sprintf("unexpected \"%s\" (=%d) expression (%T)", e.Kind(), e.Kind(), e))
	}

	// TODO: remove this return
	return nil, nil
}

func (s *Scope) evalConstSymbolExp(sym ast.Symbol) (Exp, diag.Error) {
	name := sym.Name
	pin := sym.Pin

	symbol := s.Lookup(name)
	if s == nil {
		return nil, &diag.SimpleMessageError{
			Pin:  pin,
			Text: fmt.Sprintf("type name \"%s\" refers to undefined symbol", name),
		}
	}
	if symbol.Kind != smk.Const {
		return nil, &diag.SimpleMessageError{
			Pin:  pin,
			Text: fmt.Sprintf("name \"%s\" refers to %s, not a constant", name, s.Kind),
		}
	}

	return symbol.Def.(StaticValue).Exp, nil
}

func (s *Scope) evalConstBinaryExp(exp ast.Binary) (Exp, diag.Error) {
	a, err := s.EvalConstExp(exp.A)
	if err != nil {
		return nil, err
	}
	b, err := s.EvalConstExp(exp.B)
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
			return addIntegers(s.Types, a.(Integer), b.(Integer)), nil
		case bok.Sub:
			return subIntegers(s.Types, a.(Integer), b.(Integer)), nil
		default:
			panic(fmt.Sprintf("unexpected \"%s\" (=%d) binary operator", exp.Op.Kind, exp.Op.Kind))
		}
	default:
		panic(fmt.Sprintf("unexpected \"%s\" (=%d) type (%T)", typ.Kind, typ.Kind, typ.Def))
	}
}

func addIntegers(x *TypeIndex, a, b Integer) Integer {
	switch {
	case a.Neg && b.Neg:
		return x.MakeNegInteger(a.Pin, a.Val+b.Val)
	case a.Neg && !b.Neg:
		if b.Val > a.Val {
			return x.MakeNegInteger(a.Pin, b.Val-a.Val)
		}
		return x.MakeInteger(a.Pin, a.Val-b.Val)
	case !a.Neg && b.Neg:
		if a.Val > b.Val {
			return x.MakeNegInteger(a.Pin, a.Val-b.Val)
		}
		return x.MakeInteger(a.Pin, b.Val-a.Val)
	case !a.Neg && !b.Neg:
		return x.MakeInteger(a.Pin, a.Val+b.Val)
	default:
		panic("impossible condition")
	}
}

func subIntegers(x *TypeIndex, a, b Integer) Integer {
	switch {
	case a.Neg && b.Neg:
		if a.Val > b.Val {
			return x.MakeNegInteger(a.Pin, a.Val-b.Val)
		}
		return x.MakeInteger(a.Pin, b.Val-a.Val)
	case a.Neg && !b.Neg:
		return x.MakeNegInteger(a.Pin, a.Val+b.Val)
	case !a.Neg && b.Neg:
		return x.MakeInteger(a.Pin, a.Val+b.Val)
	case !a.Neg && !b.Neg:
		if b.Val > a.Val {
			return x.MakeNegInteger(a.Pin, b.Val-a.Val)
		}
		return x.MakeInteger(a.Pin, a.Val-b.Val)
	default:
		panic("impossible condition")
	}
}

func expectInteger(exp Exp) (Integer, diag.Error) {
	n, ok := exp.(Integer)
	if ok {
		return n, nil
	}

	return Integer{}, &diag.SimpleMessageError{
		Pin:  exp.Span().Pin,
		Text: fmt.Sprintf("expected integer expression, got %s", exp.Type()),
	}
}
