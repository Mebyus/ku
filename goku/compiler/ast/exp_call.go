package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/exk"
	"github.com/mebyus/ku/goku/compiler/sm"
)

// Call represents an expression of calling something.
//
//	Call -> Chain "(" Args ")"
//	Args -> { Exp "," } // trailing comma is optional
type Call struct {
	nodeOperand

	Chain Chain
	Args  []Exp
}

var _ Exp = Call{}

func (Call) Kind() exk.Kind {
	return exk.Call
}

func (c Call) Span() sm.Span {
	return c.Chain.TailSpan()
}

func (c Call) String() string {
	var g Printer
	g.Call(c)
	return g.Output()
}
