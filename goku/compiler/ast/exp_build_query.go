package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/exk"
	"github.com/mebyus/ku/goku/compiler/srcmap"
)

// BuildQuery represents compile-time query of build variable.
//
// FormalDefinition:
//
//	BuildQuery => "#build" { "." word }
type BuildQuery struct {
	nodeOperand

	Name string
	Pin  srcmap.Pin
}

// Explicit interface implementation check.
var _ Operand = BuildQuery{}

func (BuildQuery) Kind() exk.Kind {
	return exk.BuildQuery
}

func (q BuildQuery) Span() srcmap.Span {
	return srcmap.Span{Pin: q.Pin, Len: uint32(len(q.Name))}
}

func (q BuildQuery) String() string {
	var g Printer
	g.BuildQuery(q)
	return g.Output()
}
