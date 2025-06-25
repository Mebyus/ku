package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/exk"
	"github.com/mebyus/ku/goku/compiler/srcmap"
)

type Ref struct {
	nodeOperand

	Chain Chain
}

// Explicit interface implementation check.
var _ Operand = Ref{}

func (Ref) Kind() exk.Kind {
	return exk.Ref
}

func (r Ref) Span() srcmap.Span {
	return r.Chain.TailSpan()
}

func (r Ref) String() string {
	var g Printer
	g.Ref(r)
	return g.Output()
}
