package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/exk"
	"github.com/mebyus/ku/goku/compiler/source"
)

// Symbol represents a single word token usage inside an expression.
//
// Formal definition:
//
//	Symbol => word
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

func (s Symbol) Span() source.Span {
	return source.Span{Pin: s.Pin, Len: uint32(len(s.Name))}
}

func (s Symbol) String() string {
	return s.Name
}

// DotName represents contextual symbol usage inside an expression.
//
// Formal definition:
//
//	DotName => "." word
type DotName struct {
	nodeOperand

	// Symbol name. Single word.
	Name string

	Pin source.Pin
}

// Explicit interface implementation check.
var _ Operand = Symbol{}

func (DotName) Kind() exk.Kind {
	return exk.DotName
}

func (d DotName) Span() source.Span {
	return source.Span{Pin: d.Pin, Len: uint32(len(d.Name)) + 1}
}

func (d DotName) String() string {
	return "." + d.Name
}
