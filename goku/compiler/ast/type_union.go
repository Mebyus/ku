package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/tsk"
	"github.com/mebyus/ku/goku/compiler/srcmap"
)

// Union represents union type specifier.
//
// Formal definition:
//
//	Union => "union" "{" { Field "," } "}"
type Union struct {
	// Can be nil (if struct does not have fields).
	Fields []Field

	Pin srcmap.Pin
}

var _ TypeSpec = Union{}

func (Union) Kind() tsk.Kind {
	return tsk.Union
}

func (u Union) Span() srcmap.Span {
	return srcmap.Span{Pin: u.Pin}
}

func (u Union) String() string {
	var g Printer
	g.Union(u)
	return g.Output()
}
