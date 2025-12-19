package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/exk"
	"github.com/mebyus/ku/goku/compiler/sm"
)

type Index struct {
	nodePart

	// Expression inside index.
	Exp Exp
}

// Explicit interface implementation check.
var _ Part = Index{}

func (Index) Kind() exk.Kind {
	return exk.Index
}

func (x Index) Span() sm.Span {
	return x.Exp.Span()
}

func (x Index) String() string {
	var g Printer
	g.Index(x)
	return g.Output()
}

type DerefIndex struct {
	nodePart

	// Expression inside index.
	Exp Exp
}

// Explicit interface implementation check.
var _ Part = DerefIndex{}

func (DerefIndex) Kind() exk.Kind {
	return exk.DerefIndex
}

func (x DerefIndex) Span() sm.Span {
	return x.Exp.Span()
}

func (x DerefIndex) String() string {
	var g Printer
	g.DerefIndex(x)
	return g.Output()
}
