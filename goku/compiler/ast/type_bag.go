package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/tsk"
	"github.com/mebyus/ku/goku/compiler/sm"
)

type Bag struct {
	Funs []BagFun

	Pin sm.Pin
}

var _ TypeSpec = Bag{}

func (Bag) Kind() tsk.Kind {
	return tsk.Bag
}

func (b Bag) Span() sm.Span {
	return sm.Span{Pin: b.Pin}
}

func (b Bag) String() string {
	var g Printer
	g.Bag(b)
	return g.Output()
}

// BagFun represents a single bag function inside a bag.
//
// Formal definition:
//
//	Function => Name Signature
//	Name     => word
type BagFun struct {
	Signature Signature

	Name Word
}
