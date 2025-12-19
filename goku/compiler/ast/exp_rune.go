package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/exk"
	"github.com/mebyus/ku/goku/compiler/sm"
)

type Rune struct {
	nodeOperand

	// Rune literal value represented by token.
	Val uint64

	Pin sm.Pin
}

// Explicit interface implementation check.
var _ Operand = Rune{}

func (Rune) Kind() exk.Kind {
	return exk.Rune
}

func (r Rune) Span() sm.Span {
	return sm.Span{Pin: r.Pin}
}

func (r Rune) String() string {
	var g Printer
	g.Rune(r)
	return g.Output()
}
