package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/exk"
	"github.com/mebyus/ku/goku/compiler/source"
)

// Paren represents an expression which consists of another expression and
// parenthesis surrounding it.
//
// Formal definition:
//
//	Paren => "(" Exp ")"
type Paren struct {
	nodeOperand

	// Expression surrounded by parenthesis.
	Exp Exp
}

// Explicit interface implementation check.
var _ Operand = Paren{}

func (Paren) Kind() exk.Kind {
	return exk.Paren
}

func (p Paren) Span() source.Span {
	return p.Exp.Span()
}

func (p Paren) String() string {
	var g Printer
	g.Paren(p)
	return g.Output()
}
