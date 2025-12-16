package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/tsk"
	"github.com/mebyus/ku/goku/compiler/srcmap"
)

// Chunk represents chunk type specifier.
//
// Formal definition:
//
//	Chunk => "[]" TypeSpec
type Chunk struct { // TODO: rename this to Span across project
	// Chunk element type specifier.
	Type TypeSpec
}

var _ TypeSpec = Chunk{}

func (Chunk) Kind() tsk.Kind {
	return tsk.Chunk
}

func (c Chunk) Span() srcmap.Span {
	return c.Type.Span()
}

func (c Chunk) String() string {
	var g Printer
	g.Chunk(c)
	return g.Output()
}
