package ast

import "github.com/mebyus/ku/internal/ku/sx"

// SymExp represents expression operand formed by usage of symbol name and alter zero.
//
// Formal definition:
//
//	SymZeroExp -> word ".{}"
type SymZeroExp struct {
	operand

	Name string

	Pin sx.Pin
}

// Explicit interface implementation check.
var _ Operand = &SymZeroExp{}
