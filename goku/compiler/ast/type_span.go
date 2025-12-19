package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/tsk"
	"github.com/mebyus/ku/goku/compiler/srcmap"
)

// Span represents chunk type specifier.
//
// Formal definition:
//
//	Span => "[]" TypeSpec
type Span struct { // TODO: rename this to Span across project
	// Chunk element type specifier.
	Type TypeSpec
}

var _ TypeSpec = Span{}

func (Span) Kind() tsk.Kind {
	return tsk.Span
}

func (c Span) Span() srcmap.Span {
	return c.Type.Span()
}

func (c Span) String() string {
	var g Printer
	g.Span(c)
	return g.Output()
}
