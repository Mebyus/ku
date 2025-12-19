package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/exk"
	"github.com/mebyus/ku/goku/compiler/sm"
)

// True represents usage of "true" as expression.
type True struct {
	nodeOperand

	Pin sm.Pin
}

// Explicit interface implementation check.
var _ Operand = True{}

func (True) Kind() exk.Kind {
	return exk.True
}

func (t True) Span() sm.Span {
	return sm.Span{Pin: t.Pin, Len: 4}
}

func (True) String() string {
	return "true"
}
