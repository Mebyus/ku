package ast

import "github.com/mebyus/ku/internal/ku/sx"

// While represents a loop with condition.
//
// Formal definition:
//
//	While -> "for" Exp Block
type While struct {
	stm

	Body Block

	// Loop condition. Always not nil.
	Exp Exp

	Pin sx.Pin
}

var _ Statement = &While{}
