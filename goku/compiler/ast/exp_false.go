package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/exk"
	"github.com/mebyus/ku/goku/compiler/srcmap"
)

// False represents usage of "false" as expression.
type False struct {
	nodeOperand

	Pin srcmap.Pin
}

// Explicit interface implementation check.
var _ Operand = False{}

func (False) Kind() exk.Kind {
	return exk.False
}

func (f False) Span() srcmap.Span {
	return srcmap.Span{Pin: f.Pin, Len: 5}
}

func (False) String() string {
	return "false"
}
