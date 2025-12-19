package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/exk"
	"github.com/mebyus/ku/goku/compiler/sm"
)

// False represents usage of "false" as expression.
type False struct {
	nodeOperand

	Pin sm.Pin
}

// Explicit interface implementation check.
var _ Operand = False{}

func (False) Kind() exk.Kind {
	return exk.False
}

func (f False) Span() sm.Span {
	return sm.Span{Pin: f.Pin, Len: 5}
}

func (False) String() string {
	return "false"
}
