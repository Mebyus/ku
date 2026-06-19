package ast

import "github.com/mebyus/ku/internal/ku/sx"

// TypeSpec node that represents type specifier of any kind.
type TypeSpec interface {
	_spec()
}

// Embed this to quickly implement _spec() discriminator from TypeSpec interface.
// Do not use it for anything else.
type spec struct{}

func (spec) _spec() {}

// Explicit interface implementation check.
var _ TypeSpec = spec{}

// Type represents top level custom type definition construct.
//
// Formal definitino:
//
//	Type    -> "type" Name [ "in" BagList ] TypeSpec [ ";" ]
//	BagList -> "(" { BagName "," } ")"
//	Name    -> word
//	BagName -> word
type Type struct {
	Name string

	// Optional list of bags which this type must fit into.
	// Bags []Word

	Spec TypeSpec

	// pin of type name
	Pin sx.Pin
}
