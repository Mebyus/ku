package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/exk"
	"github.com/mebyus/ku/goku/compiler/srcmap"
)

// EnumMacro represents usage of "#enum" operator as expression.
//
// Formal definition:
//
//	EnumMacro => "#enum" "(" Name "." Entry ")"
//	Name => word
//	Entry => word
type EnumMacro struct {
	nodeOperand

	Name Word

	Entry Word
}

// Explicit interface implementation check.
var _ Operand = EnumMacro{}

func (EnumMacro) Kind() exk.Kind {
	return exk.EnumMacro
}

func (e EnumMacro) Span() srcmap.Span {
	return e.Name.Span()
}

func (e EnumMacro) String() string {
	var g Printer
	g.EnumMacro(e)
	return g.Output()
}
