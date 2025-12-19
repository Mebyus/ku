package eval

import (
	"fmt"

	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/enums/bok"
	"github.com/mebyus/ku/goku/compiler/enums/uok"
	"github.com/mebyus/ku/goku/compiler/sm"
)

type Value interface {
	_value()
}

// Embed this to quickly implement _value() discriminator from Value interface.
// Do not use it for anything else.
type nodeValue struct{}

func (nodeValue) _value() {}

type Integer struct {
	nodeValue

	Val uint64

	Pin sm.Pin

	// True for negative integers.
	Neg bool
}

type String struct {
	nodeValue

	Val string
	Pin sm.Pin
}

type Boolean struct {
	nodeValue

	Pin sm.Pin
	Val bool
}

// EvalExp evaluate static expression based on given environment.
// Returns error if expression contains non-static elements.
func EvalExp(env *Env, exp ast.Exp) (Value, diag.Error) {
	switch e := exp.(type) {
	case nil:
		panic("nil expression")
	case ast.Paren:
		return EvalExp(env, e.Exp)
	case ast.BuildQuery:
		return evalBuildQuery(env, e)
	case ast.EnvQuery:
		return evalEnvQuery(env, e)
	case ast.Unary:
		return evalUnary(env, e)
	case ast.Binary:
		return evalBinary(env, e)
	case ast.Integer:
		return Integer{
			Val: e.Val,
			Pin: e.Pin,
		}, nil
	case ast.String:
		return String{
			Val: e.Val,
			Pin: e.Pin,
		}, nil
	case ast.True:
		return Boolean{
			Val: true,
			Pin: e.Pin,
		}, nil
	case ast.False:
		return Boolean{
			Val: false,
			Pin: e.Pin,
		}, nil
	default:
		return nil, &diag.SimpleMessageError{
			Pin:  exp.Span().Pin,
			Text: fmt.Sprintf("unexpected \"%s\" (=%d) expression (%T)", e.Kind(), e.Kind(), e),
		}
	}
}

func evalBuildQuery(env *Env, q ast.BuildQuery) (Value, diag.Error) {
	switch q.Name {
	case "kind":
		return Integer{
			Val: uint64(env.BuildKind),
			Pin: q.Pin,
		}, nil
	case "mode":
		return Integer{
			Val: uint64(env.BuildMode),
			Pin: q.Pin,
		}, nil
	default:
		return nil, &diag.SimpleMessageError{
			Pin:  q.Pin,
			Text: fmt.Sprintf("unknown build variable \"%s\"", q.Name),
		}
	}
}

func evalEnvQuery(env *Env, q ast.EnvQuery) (Value, diag.Error) {
	v, ok := env.m[q.Name]
	if ok {
		return v, nil
	}

	return nil, &diag.SimpleMessageError{
		Pin:  q.Pin,
		Text: fmt.Sprintf("unknown env \"%s\"", q.Name),
	}
}

func evalUnary(env *Env, u ast.Unary) (Value, diag.Error) {
	value, err := EvalExp(env, u.Exp)
	if err != nil {
		return nil, err
	}

	switch v := value.(type) {
	case Integer:
		switch u.Op.Kind {
		case uok.Plus:
			return Integer{
				Val: v.Val,
				Neg: v.Neg,
				Pin: u.Op.Pin,
			}, nil
		case uok.Minus:
			if v.Val == 0 {
				return Integer{
					Val: 0,
					Pin: u.Op.Pin,
				}, nil
			}
			return Integer{
				Val: v.Val,
				Neg: !v.Neg,
				Pin: u.Op.Pin,
			}, nil
		case uok.BitNot:
			if v.Neg {
				return nil, &diag.SimpleMessageError{
					Pin:  u.Op.Pin,
					Text: fmt.Sprintf("unary operation \"%s\" not supported on negative integer value", u.Op.Kind),
				}
			}

			return Integer{
				Val: ^v.Val,
				Pin: u.Op.Pin,
			}, nil
		default:
			return nil, &diag.SimpleMessageError{
				Pin:  u.Op.Pin,
				Text: fmt.Sprintf("unary operation \"%s\" not supported on integer value", u.Op.Kind),
			}
		}
	case String:
		return nil, &diag.SimpleMessageError{
			Pin:  u.Op.Pin,
			Text: fmt.Sprintf("unary operation \"%s\" not supported on string value", u.Op.Kind),
		}
	case Boolean:
		if u.Op.Kind == uok.Not {
			return Boolean{
				Val: !v.Val,
				Pin: u.Op.Pin,
			}, nil
		}

		return nil, &diag.SimpleMessageError{
			Pin:  u.Op.Pin,
			Text: fmt.Sprintf("unary operation \"%s\" not supported on boolean value", u.Op.Kind),
		}
	default:
		panic(fmt.Sprintf("unexpected value (%T)", v))
	}
}

func evalBinary(env *Env, bin ast.Binary) (Value, diag.Error) {
	a, err := EvalExp(env, bin.A)
	if err != nil {
		return nil, err
	}
	b, err := EvalExp(env, bin.B)
	if err != nil {
		return nil, err
	}

	switch a := a.(type) {
	case Integer:
		switch b := b.(type) {
		case Integer:
			switch bin.Op.Kind {
			case bok.Add:
				return Integer{
					Pin: a.Pin,
				}, nil
			case bok.Sub:
				return Integer{
					Pin: a.Pin,
				}, nil
			case bok.Equal:
				return Boolean{
					Val: a.Neg == b.Neg && a.Val == b.Val,
					Pin: a.Pin,
				}, nil
			case bok.NotEqual:
				return Boolean{
					Val: a.Neg != b.Neg || a.Val != b.Val,
					Pin: a.Pin,
				}, nil
			case bok.Less:
				var val bool
				switch {
				case !a.Neg && !b.Neg:
					val = a.Val < b.Val
				case a.Neg && b.Neg:
					val = a.Val > b.Val
				case a.Neg && !b.Neg:
					val = true
				case !a.Neg && b.Neg:
					val = false
				default:
					panic("unreachable")
				}
				return Boolean{
					Val: val,
					Pin: a.Pin,
				}, nil
			case bok.LessOrEqual:
				var val bool
				switch {
				case !a.Neg && !b.Neg:
					val = a.Val <= b.Val
				case a.Neg && b.Neg:
					val = a.Val >= b.Val
				case a.Neg && !b.Neg:
					val = true
				case !a.Neg && b.Neg:
					val = false
				default:
					panic("unreachable")
				}
				return Boolean{
					Val: val,
					Pin: a.Pin,
				}, nil
			case bok.Greater:
				return Boolean{
					Pin: a.Pin,
				}, nil
			case bok.GreaterOrEqual:
				return Boolean{
					Pin: a.Pin,
				}, nil
			default:
				return nil, &diag.SimpleMessageError{
					Pin:  bin.Op.Pin,
					Text: fmt.Sprintf("binary operation \"%s\" not supported on integer values", bin.Op.Kind),
				}
			}
		default:
			return nil, &diag.SimpleMessageError{
				Pin:  bin.Op.Pin,
				Text: fmt.Sprintf("binary operation \"%s\" on mismatched values (integer, %T)", bin.Op.Kind, b),
			}
		}
	case String:
		switch b := b.(type) {
		case String:
			switch bin.Op.Kind {
			case bok.Add:
				return String{
					Val: a.Val + b.Val,
					Pin: a.Pin,
				}, nil
			case bok.Equal:
				return Boolean{
					Val: a.Val == b.Val,
					Pin: a.Pin,
				}, nil
			case bok.NotEqual:
				return Boolean{
					Val: a.Val != b.Val,
					Pin: a.Pin,
				}, nil
			default:
				return nil, &diag.SimpleMessageError{
					Pin:  bin.Op.Pin,
					Text: fmt.Sprintf("binary operation \"%s\" not supported on string values", bin.Op.Kind),
				}
			}
		default:
			return nil, &diag.SimpleMessageError{
				Pin:  bin.Op.Pin,
				Text: fmt.Sprintf("binary operation \"%s\" on mismatched values (string, %T)", bin.Op.Kind, b),
			}
		}
	case Boolean:
		panic("not implemented")
	default:
		panic(fmt.Sprintf("unexpected value (%T)", a))
	}
}
