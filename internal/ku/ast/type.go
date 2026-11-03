package ast

import (
	"github.com/mebyus/ku/internal/fbuf"
	"github.com/mebyus/ku/internal/ku/sx"
)

// TypeSpec node that represents type specifier of any kind.
type TypeSpec interface {
	// Discriminator method for interface implementations.
	// Only serves as a trick to enhance Go typechecking in
	// type assertions.
	//
	// Does nothing when called.
	_spec()

	// Format type specifier as text that should be very close
	// or identical to what it was parsed from.
	//
	// Used mainly for tests and debugging.
	Format(*fbuf.Buffer)
}

// Embed this to quickly implement _spec() discriminator from TypeSpec interface.
// Do not use it for anything else.
type spec struct{}

func (spec) _spec() {}

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
