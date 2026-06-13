package ast

type If struct {
	stm

	// condition
	Exp Exp

	// true branch
	Body Block

	// can be nil if statement does not have else branch
	Else *Block
}

var _ Statement = &If{}
