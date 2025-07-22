package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/exk"
	"github.com/mebyus/ku/goku/compiler/srcmap"
)

// CheckFlag represents a macro for checking a bit flag in expression.
//
// Formal definition:
//
//	CheckFlag => "#check" "(" Exp "," Exp ")"
type CheckFlag struct {
	nodeOperand

	Exp  Exp
	Flag Exp
}

// Explicit interface implementation check.
var _ Operand = CheckFlag{}

func (CheckFlag) Kind() exk.Kind {
	return exk.CheckFlag
}

func (c CheckFlag) Span() srcmap.Span {
	return c.Exp.Span()
}

func (c CheckFlag) String() string {
	var g Printer
	g.CheckFlag(c)
	return g.Output()
}
