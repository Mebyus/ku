package ast

// Formal definition:
//
//	Fun  => "fun" Name Signature Body
//	Body => Block
//	Name => word
type Fun struct {
	Signature Signature

	Body Block
	Name Word

	Traits
}

// Formal definition:
//
//	Fun  => "test" Name Body
//	Body => Block
//	Name => word
type TestFun struct {
	Body Block
	Name Word
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
//	Param => Name ":" TypeSpec
//	Name  => word
type Param struct {
	Name Word
	Type TypeSpec
}
