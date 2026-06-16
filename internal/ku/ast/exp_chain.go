package ast

import "github.com/mebyus/ku/internal/ku/sx"

// Chain represents an operand which starts with identifier and has at least
// one chain part operator attached to it (select, call, deref, index, etc.)
type Chain struct {
	operand

	// name of identifier which starts the chain
	Name string

	// Always has at least one element.
	Parts []Part

	// pin of start identifier
	Pin sx.Pin
}

var _ Operand = &Chain{}

type Part interface {
	// Discriminator method for interface implementations.
	// Only serves as a trick to enhance Go typechecking in
	// type assertions.
	//
	// Does nothing when called.
	_part()
}

// Embed this to quickly implement _part() discriminator from Part interface.
// Do not use it for anything else.
type part struct{}

func (part) _part() {}

type Select struct {
	part

	Name string

	Pin sx.Pin
}

// Explicit interface implementation check.
var _ Part = &Select{}
