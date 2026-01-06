package stg

import (
	"fmt"

	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/enums/bgk"
	"github.com/mebyus/ku/goku/compiler/enums/tpk"
	"github.com/mebyus/ku/goku/compiler/sm"
)

func (s *Scope) translateBgenFunCall(se *SymExp, args []ast.Exp) (Exp, diag.Error) {
	k := bgk.Kind(se.Symbol.Aux)
	switch k {
	case bgk.Min:
		return s.translateMinCall(se.Pin, args)
	default:
		panic(fmt.Sprintf("unexpected \"%s\" (=%d) builtin generic function", se.Symbol.Name, k))
	}
}

func (s *Scope) translateMinCall(pin sm.Pin, args []ast.Exp) (Exp, diag.Error) {
	if len(args) == 0 {
		return nil, &diag.SimpleMessageError{
			Pin:  pin,
			Text: "min function cannot be called with 0 args",
		}
	}

	rargs := make([]Exp, 0, len(args)) // runtime args
	var cargs []Exp                    // compile-time args

	for _, arg := range args {
		a, err := s.TranslateExp(&Hint{}, arg)
		if err != nil {
			return nil, err
		}

		typ := a.Type()
		if typ.Kind != tpk.Integer {
			return nil, &diag.SimpleMessageError{
				Pin:  a.Span().Pin,
				Text: fmt.Sprintf("call to min function contains argument with %s type, which is not an integer type", typ),
			}
		}
		if typ.IsStatic() {
			cargs = append(cargs, a)
		} else {
			rargs = append(rargs, a)
		}
	}

	if len(rargs) == 0 {
		// no runtime args, calculate result at compile-time
		return s.Types.calcMin(cargs), nil
	}

	typ := rargs[0].Type()
	for _, arg := range rargs[1:] {
		t := arg.Type()
		if t != typ {
			return nil, &diag.SimpleMessageError{
				Pin:  arg.Span().Pin,
				Text: fmt.Sprintf("mismatched argument types %s and %s in call to min function", typ, t),
			}
		}
	}

	if len(cargs) != 0 {
		// simplify compile-time arguments
		// if there are more than one such argument, we can select min among them
		// and leave only one compile-time argument
		a := s.Types.calcMin(cargs)
		v := a.(*Integer)

		if !typ.IsSigned() && v.Neg {
			return nil, &diag.SimpleMessageError{
				Pin:  v.Pin,
				Text: fmt.Sprintf("negative static integer mixed with unsigned integer type %s in call to min", typ),
			}
		}

		rargs = append(rargs, a)
	}

	symbol := s.Gens.getMinInstance(UniformParamsSpec{
		Type: typ,
		Num:  uint(len(rargs)),
	})

	return &Call{
		typ:    typ,
		Symbol: symbol,
		Args:   rargs,
		Pin:    pin,
	}, nil
}

func (x *TypeIndex) calcMin(list []Exp) Exp {
	if len(list) == 1 {
		return list[0]
	}

	m := list[0].(*Integer)
	for _, e := range list[1:] {
		v := e.(*Integer)

		switch {
		case m.Neg && v.Neg:
			if v.Val > m.Val {
				m.Val = v.Val
			}
		case !m.Neg && v.Neg:
			m.Val = v.Val
			m.Neg = true
		case m.Neg && !v.Neg:
			// do nothing
		case !m.Neg && !v.Neg:
			if v.Val < m.Val {
				m.Val = v.Val
			}
		default:
			panic("impossible condition")
		}
	}
	return m
}
