package symk

// Kind indicates symbol kind.
type Kind uint8

const (
	// Zero value of Kind. Should not be used explicitly.
	//
	// Mostly a trick to detect places where Kind is left unspecified.
	// empty Kind = iota

	// Function (produced from declaration or definition).
	Fun Kind = iota

	// // Method (produced from declaration or definition).
	// Method

	// // Unit test.
	// Test

	// Custom or builtin type definition.
	Type

	// Immutable value definition (name + type + value).
	//
	// Must be compile-time constant.
	Const

	// Carrier variable, which can only be reassigned with special
	// language construct.
	Car

	// Mutable variable which can be mutated freely.
	Var

	// // Loop variable.
	// Loop

	// Alias

	// // Builtin generic function.
	// BgenFun

	// BgenFunInst

	// // Map method (before monomorphization).
	// MapMethod

	// Gen

	// Runtime function or method parameter.
	Param

	// // Method receiver.
	// Receiver

	// // Parameter which name is omitted. Such symbols are used in cases:
	// //
	// //	- function forward or external declaration
	// //	- underscore "_" parameter in function
	// //	- parameters of builtin (global) functions
	// //
	// // Such symbols always have empty Symbol.Name field and are not
	// // added to scope.
	// OmitParam

	// // Type, function or method parameter for which argument value must be known at
	// // buildtime
	// StaticParam

	// // Symbol created by importing other unit.
	// Import
)

var text = [...]string{
	// empty: "<nil>",

	Fun: "fun",
	// Method: "method",
	// Test:   "test",
	Type:  "type",
	Const: "const",
	// Alias:  "alias",
	Car: "car",
	Var: "var",
	// Gen:    "gen",
	// Import: "import",
	Param: "param",

	// Receiver:  "receiver",
	// OmitParam: "param.omit",
	// Loop:      "var.loop",

	// BgenFun:     "builtin generic function",
	// BgenFunInst: "builtin generic function instance",
	// MapMethod:   "map method",
}

func (k Kind) String() string {
	return text[k]
}
