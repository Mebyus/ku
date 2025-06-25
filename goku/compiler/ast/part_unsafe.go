package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/exk"
	"github.com/mebyus/ku/goku/compiler/srcmap"
)

// Unsafe represents chain part which selects a name from unsafe scope.
//
// Formal definition:
//
//	Unsafe => "." "unsafe" "." word
type Unsafe struct {
	nodePart

	Name string

	Pin srcmap.Pin
}

// Explicit interface implementation check.
var _ Part = Unsafe{}

func (Unsafe) Kind() exk.Kind {
	return exk.Deref
}

func (u Unsafe) Span() srcmap.Span {
	return srcmap.Span{Pin: u.Pin, Len: uint32(len(u.Name)) + 8}
}

func (u Unsafe) String() string {
	var g Printer
	g.Unsafe(u)
	return g.Output()
}
