package stg

import (
	"fmt"

	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/enums/tpk"
	"github.com/mebyus/ku/goku/compiler/sm"
)

func CheckAssign(want *Type, exp Exp) diag.Error {
	typ := exp.Type()
	if typ == want {
		return nil
	}

	if typ.Kind != want.Kind {
		switch want.Kind {
		case tpk.Custom:
			c := want.Def.(*Custom).Type
			if typ.IsStatic() && typ.Size == 0 && typ.Kind == c.Kind {
				return nil
			}
			if c.Kind == tpk.Enum {
				if c == typ {
					return nil
				}
			}
		case tpk.Integer:
			if typ.Kind == tpk.Rune {
				// TODO: check size for static runes
				return nil
			}
		}

		return &diag.SimpleMessageError{
			Pin:  exp.Span().Pin,
			Text: fmt.Sprintf("incompatible types %s and %s", want, typ),
		}
	}

	if typ.Kind == tpk.Integer && want.Kind == tpk.Integer {
		if typ.IsSigned() == want.IsSigned() && want.Size > typ.Size {
			return nil
		}
	}

	if typ.IsStatic() && typ.Size == 0 {
		return nil
	}

	return &diag.SimpleMessageError{
		Pin:  exp.Span().Pin,
		Text: fmt.Sprintf("incompatible types %s and %s", want, typ),
	}
}

func CheckCall(sig *Signature, args []Exp) diag.Error {
	if sig.Receiver != nil {
		// First element in args must be receiver expression.
		// And receiver already got checked before this.
		args = args[1:]
	}

	if len(args) != len(sig.Params) {
		var pin sm.Pin
		if len(args) != 0 {
			// TODO: how to set pin properly if there are no args?
			pin = args[0].Span().Pin
		}
		return &diag.SimpleMessageError{
			Pin:  pin,
			Text: fmt.Sprintf("call requires %d argument(s), but got %d", len(sig.Params), len(args)),
		}
	}

	for i := range len(args) {
		arg := args[i]
		param := sig.Params[i]

		err := CheckAssign(param, arg)
		if err != nil {
			return err
		}
	}

	return nil
}

func (x *TypeIndex) CheckCast(want *Type, exp Exp) diag.Error {
	typ := exp.Type()
	if typ.Kind == tpk.Integer && want.Kind == tpk.Integer {
		return nil
	}

	if want.Kind == tpk.Integer && typ.Kind == tpk.Custom && typ.Def.(*Custom).Type.Kind == tpk.Enum {
		return nil
	}

	if want.Kind == tpk.Custom && want.Def.(*Custom).Type.Kind == tpk.Integer && typ.Kind == tpk.Integer {
		return nil
	}

	if want == x.Known.Str && typ == x.Known.SpanU8 {
		return nil
	}

	return &diag.SimpleMessageError{
		Pin:  exp.Span().Pin,
		Text: fmt.Sprintf("value of %s type cannot be cast into %s", typ, want),
	}
}
