package ast

import "github.com/mebyus/ku/internal/ku/sx"

// IfBlock represents "if-then" branching statement which
// does not have else block.
type IfBlock struct {
	stm

	// condition
	Exp Exp

	// then branch
	Body Block

	// "if" pin
	Pin sx.Pin
}

var _ Statement = &IfBlock{}

// Branch represents "if-then-else" branching statement.
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

// IfFlow represents short (one line) form of "if" branching statement.
// It supports only restricted set of statements which affect control flow:
//
//   - return
//   - break
//   - gonext
//   - goto
type IfFlow struct {
	stm

	// condition
	Exp Exp

	// true branch
	Then Statement

	Pin sx.Pin
}

var _ Statement = &IfFlow{}
