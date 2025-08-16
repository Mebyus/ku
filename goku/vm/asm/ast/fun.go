package ast

import (
	"github.com/mebyus/ku/goku/compiler/srcmap"
	"github.com/mebyus/ku/goku/vm/opc"
)

type Fun struct {
	Atoms []Atom

	// List of labels placed inside function body.
	// Stored in placement order.
	Labels []string

	Name string

	Pin srcmap.Pin
}

// Label represents label name usage.
type Label struct {
	nodeOperand

	Name string
	Pin  srcmap.Pin
}

// Place represents label placement operation in function body.
type Place struct {
	nodeAtom

	Name string
	Pin  srcmap.Pin
}

// Symbol represents symbol usage operand (inside instruction).
type Symbol struct {
	nodeOperand

	Name string
	Pin  srcmap.Pin
}

// Integer represents integer usage operand (inside instruction).
type Integer struct {
	nodeOperand

	Val uint64
	Pin srcmap.Pin
}

// Register represents register usage operand (inside instruction).
type Register struct {
	nodeOperand

	Name opc.Register
	Pin  srcmap.Pin
}

// Instruction represents arbitrary instruction in function body.
type Instruction struct {
	nodeAtom

	Operands []Operand

	// Instruction mnemonic string.
	Mnemonic string

	// Optional.
	Variant string

	// Mnemonic pin.
	Pin srcmap.Pin
}
