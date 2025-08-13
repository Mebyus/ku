package ast

import "github.com/mebyus/ku/goku/compiler/srcmap"

type Fun struct {
	Atoms []Atom

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
