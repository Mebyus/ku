package ast

import "github.com/mebyus/ku/internal/ku/sx"

// Branch represents "if-else" branching statement.
// Else branch is optional.
type Branch struct {
	stm

	// condition
	Exp Exp

	// true branch
	Body Block

	Pin sx.Pin

	// can be nil if statement does not have else branch
	Else *Block
}

var _ Statement = &Branch{}

// LineIf represents short (one line) form of "if" branching statement.
type LineIf struct {
	stm

	// condition
	Exp Exp

	// true branch
	Then Statement

	Pin sx.Pin
}

var _ Statement = &LineIf{}
