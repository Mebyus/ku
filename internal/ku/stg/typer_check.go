package stg

import (
	"fmt"

	"github.com/mebyus/ku/internal/ku/enums/bop"
	"github.com/mebyus/ku/internal/ku/enums/typk"
)

func (t *Typer) checkReturnType(rt *Type, exp Exp) {
	et := exp.Type()
	if et == rt {
		return
	}

	if rt == nil {
		if et == nil {
			// TODO: maybe should issue warning when code returns
			// expression with nothing type.
			//
			// This must only occur on returning function call, which
			// in turn returns nothing.
			return
		}

		// already issued error upon converting return statement
		// should we report another one here?
		return
	}

	if rt.IsInvalid() {
		// this was already reported earlier
		return
	}

	if et == nil {
		t.report(exp.Pin(), fmt.Sprintf("expression yield empty value where %s must be in return statement", rt))
		return
	}

	if et.IsInvalid() {
		// this was already reported earlier
		return
	}

	if rt.Kind == typk.Boolean && et.Kind == typk.Boolean {
		return
	}

	if rt.Kind == typk.Integer && et.Kind == typk.Integer {
		if et.IsStatic() {
			return
		}
		if rt.IsSigned() == et.IsSigned() && et.Size <= rt.Size {
			return
		}
	}

	t.report(exp.Pin(), fmt.Sprintf("cannot use expression of %s value as %s in return statement", et, rt))
}

// does typechecking of binary expression and sets its resulting type
func (t *Typer) checkBinExpType(exp *BinExp) {
	at := exp.A.Type()
	bt := exp.B.Type()

	if at.Kind == typk.Integer && bt.Kind == typk.Integer {
		switch exp.Op.Kind {
		case bop.Add:
			t.checkAddIntegerType(exp, at, bt)
		case bop.Mul:
			t.checkMulIntegerType(exp, at, bt)
		case bop.Div:
			t.checkDivIntegerType(exp, at, bt)
		case bop.Equal:
			t.checkEqualIntegerType(exp, at, bt)
		case bop.NotEqual:
			t.checkNotEqualIntegerType(exp, at, bt)
		default:
			panic(fmt.Sprintf("unexpected %s operator", exp.Op.Kind))
		}
		return
	}

}

func (t *Typer) checkEqualIntegerType(exp *BinExp, at, bt *Type) {
	if at == bt || at.IsStatic() || bt.IsStatic() {
		exp.typ = t.common.Types.Known.Bool
		return
	}
}

func (t *Typer) checkNotEqualIntegerType(exp *BinExp, at, bt *Type) {
	if at == bt || at.IsStatic() || bt.IsStatic() {
		exp.typ = t.common.Types.Known.Bool
		return
	}
}

func (t *Typer) checkAddIntegerType(exp *BinExp, at, bt *Type) {
	if at == bt {
		exp.typ = at
		return
	}
	if at.IsStatic() {
		exp.typ = bt
		return
	}
	if bt.IsStatic() {
		exp.typ = at
		return
	}
}

func (t *Typer) checkMulIntegerType(exp *BinExp, at, bt *Type) {
	if at == bt {
		exp.typ = at
		return
	}
	if at.IsStatic() {
		exp.typ = bt
		return
	}
	if bt.IsStatic() {
		exp.typ = at
		return
	}
}

func (t *Typer) checkDivIntegerType(exp *BinExp, at, bt *Type) {
	if at == bt {
		exp.typ = at
		return
	}
	if at.IsStatic() {
		exp.typ = bt
		return
	}
	if bt.IsStatic() {
		exp.typ = at
		return
	}
}
