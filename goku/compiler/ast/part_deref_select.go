package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/exk"
	"github.com/mebyus/ku/goku/compiler/sm"
)

type DerefSelect struct {
	nodePart

	Name Word
}

// Explicit interface implementation check.
var _ Part = DerefSelect{}

func (DerefSelect) Kind() exk.Kind {
	return exk.DerefSelect
}

func (d DerefSelect) Span() sm.Span {
	return d.Name.Span()
}

func (d DerefSelect) String() string {
	var g Printer
	g.DerefSelect(d)
	return g.Output()
}
