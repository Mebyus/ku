package stg

import "github.com/mebyus/ku/internal/ku/sx"

// Statement node that represents statement of any kind.
type Statement interface {
	// Discriminator method for interface implementations.
	// Only serves as a trick to enhance Go typechecking in
	// type assertions.
	//
	// Does nothing when called.
	_stm()
}

// Embed this to quickly implement _stm() discriminator from Statement interface.
// Do not use it for anything else.
type stm struct{}

// Explicit interface implementation check.
var _ Statement = stm{}

func (stm) _stm() {}

// Block represents block statement or function body.
type Block struct {
	stm

	Scope Scope

	Nodes []Statement

	Pin sx.Pin
}

var _ Statement = &Block{}

// Ret represents return statement.
type Return struct {
	stm

	// Can be nil, if return does not have expression.
	Exp Exp

	Pin sx.Pin
}

var _ Statement = &Return{}
