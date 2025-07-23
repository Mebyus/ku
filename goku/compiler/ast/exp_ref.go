package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/exk"
	"github.com/mebyus/ku/goku/compiler/srcmap"
)

type GetRef struct {
	nodeOperand

	Chain Chain
}

// Explicit interface implementation check.
var _ Operand = GetRef{}

func (GetRef) Kind() exk.Kind {
	return exk.Ref
}

func (r GetRef) Span() srcmap.Span {
	return r.Chain.TailSpan()
}

func (r GetRef) String() string {
	var g Printer
	g.GetRef(r)
	return g.Output()
}
