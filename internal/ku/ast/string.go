package ast

import "github.com/mebyus/ku/internal/ku/sx"

// String represents a single string token usage inside the tree.
type String struct {
	operand

	// String value represented by token. Does not contain quotes.
	// Empty string token means empty Val here.
	Val string

	Pin sx.Pin

	// Auxiliary information about the token.
	Aux uint32
}

// Explicit interface implementation check.
var _ Operand = &String{}
