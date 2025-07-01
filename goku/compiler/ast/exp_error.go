package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/exk"
	"github.com/mebyus/ku/goku/compiler/srcmap"
)

// ErrorId represents usage of "#error" operator as expression.
//
// Formal definition:
//
//	ErrorId => "#error" "(" word ")"
type ErrorId struct {
	nodeOperand

	// Name referenced by operator.
	Name Word
}

// Explicit interface implementation check.
var _ Operand = ErrorId{}

func (ErrorId) Kind() exk.Kind {
	return exk.ErrorId
}

func (e ErrorId) Span() srcmap.Span {
	return e.Name.Span()
}

func (e ErrorId) String() string {
	var g Printer
	g.ErrorId(e)
	return g.Output()
}
