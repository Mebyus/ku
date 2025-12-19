package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/tsk"
	"github.com/mebyus/ku/goku/compiler/sm"
)

// Array represents array type specifier.
//
// Formal definition:
//
//	Array => "[" Size "]" TypeSpec
//
//	Size => Exp
type Array struct {
	// May be nil if array is specified as auto-len.
	Size Exp

	// Array element type specifier.
	Type TypeSpec
}

var _ TypeSpec = Array{}

func (Array) Kind() tsk.Kind {
	return tsk.Array
}

func (a Array) Span() sm.Span {
	return a.Size.Span()
}

func (a Array) String() string {
	var g Printer
	g.Array(a)
	return g.Output()
}
