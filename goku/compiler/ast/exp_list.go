package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/exk"
	"github.com/mebyus/ku/goku/compiler/sm"
)

// List represents a literal with list of expressions.
//
//	List => "[" { Exp "," } "]" // trailing comma is optional
type List struct {
	nodeOperand

	Exps []Exp

	Pin sm.Pin
}

var _ Exp = List{}

func (List) Kind() exk.Kind {
	return exk.List
}

func (l List) Span() sm.Span {
	return sm.Span{Pin: l.Pin}
}

func (l List) String() string {
	var g Printer
	g.List(l)
	return g.Output()
}
