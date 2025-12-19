package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/exk"
	"github.com/mebyus/ku/goku/compiler/sm"
)

// ArrayLen represents a macro for compile-time calculation of array length.
//
// Formal definition:
//
//	ArrayLen => "#len" "(" Exp ")"
type ArrayLen struct {
	nodeOperand

	Exp Exp
}

// Explicit interface implementation check.
var _ Operand = ArrayLen{}

func (ArrayLen) Kind() exk.Kind {
	return exk.ArrayLen
}

func (l ArrayLen) Span() sm.Span {
	return l.Exp.Span()
}

func (l ArrayLen) String() string {
	var g Printer
	g.ArrayLen(l)
	return g.Output()
}
