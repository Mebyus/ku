package ast

import "github.com/mebyus/ku/internal/ku/sx"

// Var represents a statement that introduces a new variable
// inside current scope.
type Var struct {
	stm

	// May be nil if variable type was not specified
	// explicitly.
	Type TypeSpec

	// Expression which sets variable initial value.
	// May be nil if such expression was omitted,
	// in that case default zero value is assigned
	// implicitly.
	Exp Exp

	// Variable name.
	Name string

	// Name pin.
	Pin sx.Pin
}

var _ Statement = &Var{}

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
