package ast

import "github.com/mebyus/ku/internal/ku/sx"

// Integer represents a single integer token usage inside the tree.
type Integer struct {
	operand

	// Integer value represented by token.
	Val uint64

	Pin sx.Pin

	// Auxiliary information about the token.
	Aux uint32
}

// Explicit interface implementation check.
var _ Operand = &Integer{}
