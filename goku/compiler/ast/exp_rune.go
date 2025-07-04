package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/exk"
	"github.com/mebyus/ku/goku/compiler/srcmap"
)

type Rune struct {
	nodeOperand

	// Rune literal value represented by token.
	Val uint64

	Pin srcmap.Pin
}

// Explicit interface implementation check.
var _ Operand = Rune{}

func (Rune) Kind() exk.Kind {
	return exk.Rune
}

func (r Rune) Span() srcmap.Span {
	return srcmap.Span{Pin: r.Pin}
}

func (r Rune) String() string {
	var g Printer
	g.Rune(r)
	return g.Output()
}
