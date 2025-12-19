package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/exk"
	"github.com/mebyus/ku/goku/compiler/sm"
)

type Deref struct {
	nodePart

	Pin sm.Pin
}

// Explicit interface implementation check.
var _ Part = Deref{}

func (Deref) Kind() exk.Kind {
	return exk.Deref
}

func (d Deref) Span() sm.Span {
	return sm.Span{Pin: d.Pin}
}

func (d Deref) String() string {
	var g Printer
	g.Deref(d)
	return g.Output()
}
