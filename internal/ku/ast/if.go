package ast

import "github.com/mebyus/ku/internal/ku/sx"

type If struct {
	stm

	// condition
	Exp Exp

	// true branch
	Body Block

	Pin sx.Pin

	// can be nil if statement does not have else branch
	Else *Block
}

var _ Statement = &If{}
