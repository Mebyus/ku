package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/exk"
	"github.com/mebyus/ku/goku/compiler/source"
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

func (c Chain) TailSpan() source.Span {
	return c.Parts[len(c.Parts)-1].Span()
}

func (c Chain) Span() source.Span {
	return c.Start.Span()
}

func (c Chain) String() string {
	var g Printer
	g.Chain(c)
	return g.Output()
}
