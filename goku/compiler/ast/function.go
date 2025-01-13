package ast

// Formal definition:
//
//	Function => "fun" Name Signature Body
//	Body     => Block
//	Name     => word
type Function struct {
	Signature Signature

	Body Block
	Name Word

	Traits
}

// Formal definition:
//
//	Function => "stub" "fun" Name Signature
//	Name     => word
type FunStub struct {
	Signature Signature

	Name Word

	Traits
}

// Formal definition:
//
//	Signature => "(" ParamList ")" [ "=>" ( Result | "never" ) ]
//	ParamList => { Param "," } // last comma is optional
//	Result    => TypeSpec
type Signature struct {
	// Equals nil if there are no parameters in signature
	Params []Param

	// Equals nil if function returns nothing or never returns
	Result TypeSpec

	// Equals true if function never returns
	Never bool
}

// Param can represent a single:
//
//   - parameter in function signature
//   - field in struct definition
//
// Formal definition:
//
//	Param => Name ":" TypeSpec
//	Name  => word
type Param struct {
	Name Word
	Type TypeSpec
}
