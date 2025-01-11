package ast

import (
	"github.com/mebyus/ku/goku/enums/exk"
	"github.com/mebyus/ku/goku/source"
)

type Paren struct {
	// Expression inside parenthesis.
	Exp Exp
}

var _ Exp = Paren{}

func (Paren) Kind() exk.Kind {
	return exk.Paren
}

func (p Paren) Span() source.Span {
	return p.Exp.Span()
}

func (p Paren) String() string {
	var g Printer
	g.Paren(p)
	return g.Str()
}
