package ast

import "github.com/mebyus/ku/internal/ku/sx"

// Create represents a statement that define new (or reassigns old) fixed variable.
//
// Formal definition:
//
//	Create -> Left ":=" Value ";"
//	Left   -> word
//	Value  -> Exp
type Create struct {
	stm

	// Value which is used to perform the alteration of target.
	Exp Exp

	// Name of fixed variable.
	Name string

	// Name pin.
	Pin sx.Pin
}

var _ Statement = &Create{}
