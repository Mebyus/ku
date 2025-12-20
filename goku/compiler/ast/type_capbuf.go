package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/tsk"
	"github.com/mebyus/ku/goku/compiler/sm"
)

// Span represents chunk type specifier.
//
// Formal definition:
//
//	CapBuf => "[^]" TypeSpec
type CapBuf struct {
	// Chunk element type specifier.
	Type TypeSpec
}

var _ TypeSpec = CapBuf{}

func (CapBuf) Kind() tsk.Kind {
	return tsk.CapBuf
}

func (c CapBuf) Span() sm.Span {
	return c.Type.Span()
}

func (c CapBuf) String() string {
	var g Printer
	g.CapBuf(c)
	return g.Output()
}
