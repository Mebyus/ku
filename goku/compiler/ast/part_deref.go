package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/exk"
	"github.com/mebyus/ku/goku/compiler/srcmap"
)

type Deref struct {
	nodePart

	Pin srcmap.Pin
}

// Explicit interface implementation check.
var _ Part = Deref{}

func (Deref) Kind() exk.Kind {
	return exk.Deref
}

func (d Deref) Span() srcmap.Span {
	return srcmap.Span{Pin: d.Pin}
}

func (x Deref) String() string {
	var g Printer
	g.Deref(x)
	return g.Output()
}
