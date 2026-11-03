package ast

import (
	"github.com/mebyus/ku/internal/fbuf"
	"github.com/mebyus/ku/internal/ku/sx"
)

// SymZero represents expression operand formed by usage of symbol
// name and alter zero.
//
// Formal definition:
//
//	SymZero -> word ".{}"
type SymZero struct {
	operand

	Name string

	Pin sx.Pin
}

// Explicit interface implementation check.
var _ Operand = &SymZero{}

func (z *SymZero) Format(buf *fbuf.Buffer) {
	buf.Puts(z.Name)
	buf.Puts(".{}")
}
