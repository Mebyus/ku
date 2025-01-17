package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/exk"
	"github.com/mebyus/ku/goku/compiler/source"
)

// Symbol represents a single word token usage inside an expression.
type Symbol struct {
	nodeOperand

	// Symbol name. Single word.
	Name string

	Pin source.Pin
}

// Explicit interface implementation check.
var _ Operand = Symbol{}

func (Symbol) Kind() exk.Kind {
	return exk.Symbol
}

func (d Symbol) Span() source.Span {
	return source.Span{Pin: d.Pin, Len: uint32(len(d.Name))}
}

func (d Symbol) String() string {
	return d.Name
}
