package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/exk"
	"github.com/mebyus/ku/goku/compiler/source"
)

type Select struct {
	nodePart

	Name Word
}

// Explicit interface implementation check.
var _ Part = Select{}

func (Select) Kind() exk.Kind {
	return exk.Select
}

func (s Select) Span() source.Span {
	return s.Name.Span()
}

func (s Select) String() string {
	var g Printer
	return g.Output()
}
