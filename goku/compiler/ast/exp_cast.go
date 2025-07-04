package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/exk"
	"github.com/mebyus/ku/goku/compiler/srcmap"
)

// Cast represents a compile-time type case of expression.
//
// Formal definition:
//
//	Cast => "#cast" "(" TypeSpec "," Exp ")"
type Cast struct {
	nodeOperand

	Type TypeSpec
	Exp  Exp
}

// Explicit interface implementation check.
var _ Operand = Cast{}

func (Cast) Kind() exk.Kind {
	return exk.Cast
}

func (c Cast) Span() srcmap.Span {
	return c.Type.Span()
}

func (c Cast) String() string {
	var g Printer
	g.Cast(c)
	return g.Output()
}
