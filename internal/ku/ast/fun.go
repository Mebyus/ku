package ast

import "github.com/mebyus/ku/internal/ku/sx"

// Fun represents named function definition (signature + body).
//
// Formal definition:
//
//	Fun  -> "fun" Name Signature Body
//	Body -> Block
//	Name -> word
type Fun struct {
	Sig  Signature
	Body Block
	Name string

	// Function name pin.
	Pin sx.Pin
}

// Signature represents function signature (params + return result).
//
// Formal definition:
//
//	Signature -> "(" ParamList ")" [ "->" ( Result | "never" ) ]
//	ParamList -> { Param "," } // last comma is optional
//	Result    -> TypeSpec
type Signature struct {
	// Equals nil if there are no parameters in signature.
	Params []Param

	// Equals nil if function returns nothing or never returns.
	Result TypeSpec

	// Equals true if function never returns.
	Never bool
}

// Param represents a single parameter in function signature.
//
// Formal definition:
//
//	Param -> Name ":" TypeSpec
//	Name  -> word
type Param struct {
	Type TypeSpec
	Name string

	// Parameter name pin.
	Pin sx.Pin
}
