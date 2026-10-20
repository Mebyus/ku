package ast

import "github.com/mebyus/ku/internal/ku/sx"

// InvokeSymbol represents statement formed by function-like
// call to a symbol.
type InvokeSymbol struct {
	stm

	// can be nil if call does not have arguments
	Args []Exp

	// name of the called symbol
	Name string

	// name pin
	Pin sx.Pin
}

var _ Statement = &InvokeSymbol{}
