package ast

import "github.com/mebyus/ku/internal/ku/sx"

// AssignSymbol represents statement that assign a new value to
// existing symbol.
//
// Formal definition:
//
//	AssignSymbol -> Name = Exp ";"
//	Name         -> word
type AssignSymbol struct {
	stm

	// Always not nil
	Exp Exp

	// Name of symbol baing assigned.
	Name string

	// Name pin.
	Pin sx.Pin
}

var _ Statement = &AssignSymbol{}
