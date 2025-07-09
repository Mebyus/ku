package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/exk"
	"github.com/mebyus/ku/goku/compiler/srcmap"
)

// Tint represents a compile-time type case of expression.
//
// Formal definition:
//
//	Tint => "tint" "(" TypeSpec "," Exp ")"
type Tint struct {
	nodeOperand

	Type TypeSpec
	Exp  Exp
}

// Explicit interface implementation check.
var _ Operand = Tint{}

func (Tint) Kind() exk.Kind {
	return exk.Tint
}

func (t Tint) Span() srcmap.Span {
	return t.Type.Span()
}

func (t Tint) String() string {
	var g Printer
	g.Tint(t)
	return g.Output()
}
