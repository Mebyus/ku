package ast

import "github.com/mebyus/ku/internal/ku/sx"

// Array represents fixed size array type specifier.
//
// Formal definition:
//
//	Array -> "[" Exp "]" TypeSpec
type Array struct {
	spec

	// Always not nil.
	Exp Exp

	// Array element type specifier.
	Type TypeSpec

	// Open square pin.
	Pin sx.Pin
}

// Explicit interface implementation check.
var _ TypeSpec = &Array{}
