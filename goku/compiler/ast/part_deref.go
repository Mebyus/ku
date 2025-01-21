package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/exk"
	"github.com/mebyus/ku/goku/compiler/source"
)

type Deref struct {
	nodePart

	Pin source.Pin
}

// Explicit interface implementation check.
var _ Part = Deref{}

func (Deref) Kind() exk.Kind {
	return exk.Deref
}

func (d Deref) Span() source.Span {
	return source.Span{Pin: d.Pin}
}

func (x Deref) String() string {
	var g Printer
	g.Deref(x)
	return g.Output()
}
