package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/exk"
	"github.com/mebyus/ku/goku/compiler/sm"
)

// Float represents a single float token usage inside the tree.
type Float struct {
	nodeOperand

	// Float value represented by token.
	Val string

	Pin sm.Pin
}

// Explicit interface implementation check.
var _ Operand = Float{}

func (Float) Kind() exk.Kind {
	return exk.Float
}

func (f Float) Span() sm.Span {
	return sm.Span{Pin: f.Pin, Len: uint32(len(f.String()))}
}

func (f Float) String() string {
	return f.Val
}
