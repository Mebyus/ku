package stg

import "github.com/mebyus/ku/internal/ku/sx"

// Statement node that represents statement of any kind.
type Statement interface {
	Pin() sx.Pin

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

func (stm) _stm() {}

// Block represents block statement or function body.
type Block struct {
	stm

	Scope Scope

	Nodes []Statement

	pin sx.Pin

	// Total number of potential exits (return + never) inside this block.
	// Count includes other blocks recursively.
	Exits uint32

	ExitType ExitType
}

// Explicit interface implementation check.
var _ Statement = &Block{}

func (b *Block) Pin() sx.Pin {
	return b.pin
}

// Describes how statement (or whole code block) behaves regarding exit points execution flow.
type ExitType uint8

const (
	// in schemes:
	// o - means execution pass to next node
	// x - means exit

	// Node never exits. Means execution always passes through this node.
	//
	// Execution scheme:
	//	o -> {o} -> o
	//
	// Zero value of this type.
	ExitNever ExitType = iota

	// Node may or may not exit upon execution depending on runtime conditions.
	// Means execution may or may not pass through this node.
	//
	// Execution scheme:
	//	o -> {o} -> o
	//	o -> {x}
	ExitBranch

	// Node always exits. Means execution never passes through this node.
	//
	// Execution scheme:
	//	o -> {x}
	ExitAlways
)

// Ret represents return statement.
type Return struct {
	stm

	// Can be nil, if return does not have expression.
	Exp Exp

	pin sx.Pin
}

var _ Statement = &Return{}

func (r *Return) Pin() sx.Pin {
	return r.pin
}

type If struct {
	stm

	// true branch code block
	Body Block

	// condition, must have boolean type
	Exp Exp

	pin sx.Pin

	// can be nil, if statement does not have "else" branch
	Else *Block
}

var _ Statement = &If{}

func (f *If) Pin() sx.Pin {
	return f.pin
}
