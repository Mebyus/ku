package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/exk"
	"github.com/mebyus/ku/goku/compiler/source"
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

func (x Index) Span() source.Span {
	return x.Exp.Span()
}

func (x Index) String() string {
	var g Printer
	g.Index(x)
	return g.Output()
}
