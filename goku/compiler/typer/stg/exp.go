package stg

import (
	"fmt"

	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/enums/bok"
	"github.com/mebyus/ku/goku/compiler/enums/smk"
	"github.com/mebyus/ku/goku/compiler/enums/tpk"
	"github.com/mebyus/ku/goku/compiler/sm"
)

// Exp node that represents an arbitrary expression.
type Exp interface {
	Type() *Type

	Span() sm.Span

	// Use only for debugging.
	String() string
}

func (s *Scope) TranslateExp(exp ast.Exp) (Exp, diag.Error) {
	switch e := exp.(type) {
	case ast.Nil:
		return s.Types.MakeNil(e.Pin), nil
	case ast.Integer:
		return s.Types.MakeInteger(e.Pin, e.Val), nil
	case ast.String:
		return s.Types.MakeString(e.Pin, e.Val), nil
	case ast.Symbol:
		return s.translateSymbolExp(e)
	case ast.Binary:
		return s.translateBinaryExp(e)
	case ast.Chain:
		panic("not implemented")
	case ast.Pack:
		return s.translatePackExp(e)
	default:
		panic(fmt.Sprintf("unexpected %s (=%d) expression (%T)", e.Kind(), e.Kind(), e))
	}
}

func (s *Scope) translatePackExp(exp ast.Pack) (*Pack, diag.Error) {
	list := make([]Exp, 0, len(exp.List))
	for _, e := range exp.List {
		x, err := s.TranslateExp(e)
		if err != nil {
			return nil, err
		}
		list = append(list, x)
	}

	return s.Types.MakePack(list), nil
}

func (s *Scope) translateSymbolExp(sym ast.Symbol) (Exp, diag.Error) {
	name := sym.Name
	pin := sym.Pin

	symbol := s.Lookup(name)
	if symbol == nil {
		return nil, &diag.SimpleMessageError{
			Pin:  pin,
			Text: fmt.Sprintf("name \"%s\" refers to undefined symbol", name),
		}
	}

	switch symbol.Kind {
	case smk.Const:
		return symbol.Def.(StaticValue).Exp, nil
	case smk.Var:
		return &VarExp{Pin: pin, Symbol: symbol}, nil
	case smk.Param:
		// TODO: do we need separate type of expressions for param symbol?
		return &VarExp{Pin: pin, Symbol: symbol}, nil
	default:
		return nil, &diag.SimpleMessageError{
			Pin:  pin,
			Text: fmt.Sprintf("symbol \"%s\" refers to %s, which cannot used as operand or expression", name, symbol.Kind),
		}
	}
}

func (s *Scope) translateBinaryExp(exp ast.Binary) (Exp, diag.Error) {
	a, err := s.TranslateExp(exp.A)
	if err != nil {
		return nil, err
	}
	b, err := s.TranslateExp(exp.B)
	if err != nil {
		return nil, err
	}

	ta := a.Type()
	tb := b.Type()
	if ta.IsStatic() && tb.IsStatic() {
		return s.evalConstBinaryExp(a, b, exp.Op)
	}

	typ, err := s.Types.deduceBinaryExpType(a, b, exp.Op)
	if err != nil {
		return nil, err
	}

	return &Binary{
		typ: typ,
		Op:  exp.Op,
		A:   a,
		B:   b,
	}, nil
}

// type checks binary expression and returns its resulting type
func (x *TypeIndex) deduceBinaryExpType(a, b Exp, op BinOp) (*Type, diag.Error) {
	ta := a.Type()
	tb := b.Type()
	if ta == tb {
		// TODO: check operator
		//
		// Boths types are equal and therefore cannot be static since
		// this function is only used when at most one of a or b is static.
		return x.checkBinaryForType(a.Type(), op)
	}

	if ta.IsStatic() {
		return x.deduceBinaryExpTypeA(a, b, op)
	}

	if tb.IsStatic() {
		return x.deduceBinaryExpTypeB(a, b, op)
	}

	panic("not implemented")
}

// type checks binary expression and returns its resulting type
//
// a has static type, b has runtime type
func (x *TypeIndex) deduceBinaryExpTypeA(a, b Exp, op BinOp) (*Type, diag.Error) {
	ta := a.Type()
	tb := b.Type()

	switch tb.Kind {
	case tpk.Integer:
		if ta.Kind != tpk.Integer {
			return nil, &diag.SimpleMessageError{
				Pin:  op.Pin,
				Text: fmt.Sprintf("binary operation on incompatible types %s and %s", ta, tb),
			}
		}
		if ta.Size == 0 {
			return x.getBinaryForIntegerType(tb, op)
		}
		panic("not implemented")
	case tpk.Pointer, tpk.ArrayPointer:
		if ta.Kind != tpk.Nil {
			return nil, &diag.SimpleMessageError{
				Pin:  op.Pin,
				Text: fmt.Sprintf("binary operation on incompatible types %s and %s", ta, tb),
			}
		}

		switch op.Kind {
		case bok.Equal, bok.NotEqual:
			return x.Known.Bool, nil
		default:
			return nil, &diag.SimpleMessageError{
				Pin:  op.Pin,
				Text: fmt.Sprintf("operation %s is not defined for pointer type", op.Kind),
			}
		}
	default:
		return nil, &diag.SimpleMessageError{
			Pin:  op.Pin,
			Text: fmt.Sprintf("type %s does not have binary operations", tb),
		}
	}
}

// type checks binary expression and returns its resulting type
//
// a has runtime type, b has static type
func (x *TypeIndex) deduceBinaryExpTypeB(a, b Exp, op BinOp) (*Type, diag.Error) {
	return x.deduceBinaryExpTypeA(b, a, op)
}

// returns resulting type of binary operation on integer types
func (x *TypeIndex) getBinaryForIntegerType(typ *Type, op BinOp) (*Type, diag.Error) {
	switch op.Kind {
	case bok.Add, bok.Sub, bok.Mul, bok.Mod, bok.Div, bok.BitAnd, bok.BitOr, bok.Xor, bok.LeftShift, bok.RightShift:
		return typ, nil
	case bok.Equal, bok.NotEqual, bok.Greater, bok.GreaterOrEqual, bok.Less, bok.LessOrEqual:
		return x.Known.Bool, nil
	default:
		return nil, &diag.SimpleMessageError{
			Pin:  op.Pin,
			Text: fmt.Sprintf("operation %s is not defined for integer type", op.Kind),
		}
	}
}

// Returns resulting type for binary operation when both operands' types are equal
// and not static.
func (x *TypeIndex) checkBinaryForType(typ *Type, op BinOp) (*Type, diag.Error) {
	switch typ.Kind {
	case tpk.Integer:
		return x.getBinaryForIntegerType(typ, op)
	case tpk.String:
		switch op.Kind {
		case bok.Equal, bok.NotEqual:
			return x.Known.Bool, nil
		default:
			return nil, &diag.SimpleMessageError{
				Pin:  op.Pin,
				Text: fmt.Sprintf("operation %s is not defined for runtime string type", op.Kind),
			}
		}
	case tpk.Boolean:
		switch op.Kind {
		case bok.Equal, bok.NotEqual:
			return x.Known.Bool, nil
		default:
			return nil, &diag.SimpleMessageError{
				Pin:  op.Pin,
				Text: fmt.Sprintf("operation %s is not defined for boolean type", op.Kind),
			}
		}
	default:
		return nil, &diag.SimpleMessageError{
			Pin:  op.Pin,
			Text: fmt.Sprintf("type %s does not have binary operations", typ),
		}
	}
}
