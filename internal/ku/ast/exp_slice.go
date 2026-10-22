package ast

import "github.com/mebyus/ku/internal/ku/sx"

// SliceAll
//
// SliceAll -> Word  "[:]"
// SliceAll -> Chain "[:]"
type SliceAll struct {
	operand

	Chain Chain

	Pin sx.Pin
}

var _ Operand = &SliceAll{}
