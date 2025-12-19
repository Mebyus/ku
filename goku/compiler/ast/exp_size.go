package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/exk"
	"github.com/mebyus/ku/goku/compiler/sm"
)

// Size represents a compile-time size query of type (or expression value type).
//
// Formal definition:
//
//	Size => "#size" "(" Exp ")"
type Size struct {
	nodeOperand

	// Contains either Exp or TypeSpec.
	Exp TypeSpec // TODO: allow expressions and type specifiers here
}

// Explicit interface implementation check.
var _ Operand = Size{}

func (Size) Kind() exk.Kind {
	return exk.Size
}

func (s Size) Span() sm.Span {
	return s.Exp.Span()
}

func (s Size) String() string {
	var g Printer
	g.Size(s)
	return g.Output()
}
