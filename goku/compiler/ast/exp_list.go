package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/exk"
	"github.com/mebyus/ku/goku/compiler/srcmap"
)

// List represents a literal with list of expressions.
//
//	List => "[" { Exp "," } "]" // trailing comma is optional
type List struct {
	nodeOperand

	Exps []Exp

	Pin srcmap.Pin
}

var _ Exp = List{}

func (List) Kind() exk.Kind {
	return exk.List
}

func (l List) Span() srcmap.Span {
	return srcmap.Span{Pin: l.Pin}
}

func (l List) String() string {
	var g Printer
	g.List(l)
	return g.Output()
}
