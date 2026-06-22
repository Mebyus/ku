package ast

import "github.com/mebyus/ku/internal/ku/sx"

type Rune struct {
	operand

	// Rune literal value represented by token.
	Val uint64

	Pin sx.Pin
}

// Explicit interface implementation check.
var _ Operand = &Rune{}
