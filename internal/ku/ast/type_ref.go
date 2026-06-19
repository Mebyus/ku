package ast

import "github.com/mebyus/ku/internal/ku/sx"

// Pointer represents reference type specifier.
//
// Formal definition:
//
//	Ref -> "&" TypeSpec
type Ref struct {
	spec

	// Type to which pointer refers to.
	Type TypeSpec

	Pin sx.Pin
}

var _ TypeSpec = &Ref{}
