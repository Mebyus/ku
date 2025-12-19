package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/exk"
	"github.com/mebyus/ku/goku/compiler/sm"
)

// Nil represents usage of "nil" keyword as expression.
type Nil struct {
	nodeOperand

	Pin sm.Pin
}

// Explicit interface implementation check.
var _ Operand = Nil{}

func (Nil) Kind() exk.Kind {
	return exk.Nil
}

func (n Nil) Span() sm.Span {
	return sm.Span{Pin: n.Pin, Len: uint32(len(n.String()))}
}

func (Nil) String() string {
	return "nil"
}
