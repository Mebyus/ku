package ast

import "github.com/mebyus/ku/internal/ku/sx"

// TypeName represents type specifier in form of type name usage via
// single identifier.
//
// Formal definition:
//
//	TypeName -> word
type TypeName struct {
	spec

	Name string
	Pin  sx.Pin
}

// Explicit interface implementation check.
var _ TypeSpec = &TypeName{}
