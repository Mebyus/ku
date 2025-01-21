package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/tsk"
	"github.com/mebyus/ku/goku/compiler/source"
)

// Array represents array type specifier.
//
// Formal definition:
//
//	Array => "[" Size "]" TypeSpec
//
//	Size => Exp
type Array struct {
	Size Exp

	// Chunk element type specifier.
	Type TypeSpec
}

var _ TypeSpec = Array{}

func (Array) Kind() tsk.Kind {
	return tsk.Array
}

func (a Array) Span() source.Span {
	return a.Size.Span()
}

func (a Array) String() string {
	var g Printer
	g.Array(a)
	return g.Output()
}
