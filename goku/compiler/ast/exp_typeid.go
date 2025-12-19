package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/exk"
	"github.com/mebyus/ku/goku/compiler/sm"
)

// TypeId represents usage of "#typeid" operator as expression.
//
// Formal definition:
//
//	TypeId => "#typeid" "(" word ")"
type TypeId struct {
	nodeOperand

	// Name referenced by operator.
	Name Word
}

// Explicit interface implementation check.
var _ Operand = TypeId{}

func (TypeId) Kind() exk.Kind {
	return exk.TypeId
}

func (t TypeId) Span() sm.Span {
	return t.Name.Span()
}

func (t TypeId) String() string {
	var g Printer
	g.TypeId(t)
	return g.Output()
}
