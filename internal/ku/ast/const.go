package ast

import "github.com/mebyus/ku/internal/ku/sx"

// Const represents constant definition statement.
//
// Formal definition:
//
//	Const -> "const" Name ":" [ TypeSpec ] "=" Exp ";"
type Const struct {
	stm

	// Can be nil if constant type is not specified.
	Type TypeSpec

	// Specifies constant init value expression.
	//
	// Always not nil.
	Exp Exp

	Name string

	// Constant name pin.
	Pin sx.Pin
}

var _ Statement = &Const{}
