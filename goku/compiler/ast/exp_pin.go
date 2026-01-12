package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/exk"
	"github.com/mebyus/ku/goku/compiler/sm"
)

// PinExp represents "#pin" token used as operand in expression.
type PinExp struct {
	nodeOperand

	Pin sm.Pin
}

// Explicit interface implementation check.
var _ Operand = PinExp{}

func (PinExp) Kind() exk.Kind {
	return exk.Pin
}

func (r PinExp) Span() sm.Span {
	return sm.Span{Pin: r.Pin}
}

func (r PinExp) String() string {
	return "#pin"
}
