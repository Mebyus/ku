package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/exk"
	"github.com/mebyus/ku/goku/compiler/sm"
)

type Chain struct {
	nodeOperand

	Start Word

	// Always has at least one element.
	Parts []Part
}

var _ Operand = Chain{}

func (Chain) Kind() exk.Kind {
	return exk.Chain
}

func (c Chain) TailSpan() sm.Span {
	if len(c.Parts) == 0 {
		return c.Span()
	}
	return c.Parts[len(c.Parts)-1].Span()
}

func (c Chain) Span() sm.Span {
	return c.Start.Span()
}

func (c Chain) String() string {
	var g Printer
	g.Chain(c)
	return g.Output()
}
