package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/tsk"
	"github.com/mebyus/ku/goku/compiler/source"
)

// Tuple represents an ordered collection of type specifiers. Tuple itself is also
// a type specifier. Tuple types are used to specify return types of functions
// with multiple return values.
//
//	fun foo(a: str) => (u32, bool)
//
// There are two special cases of tuples which compiler recognizes and does some
// behind-the-scenes magic:
//
//	()             // empty tuple type
//	(ExampleType)  // tuple type with one element
//
// The former is equivalent to an empty type, The latter is transformed to be just
// an ExampleType without a tuple wrapper. Thus the following pairs of functions
// are equivalent for the purposes of type checking and following compilation
// phases:
//
//	fun foo(a: str) => ()
//	fun bar(a: u32, b: u32) => (u32)
//
//	// The above two functions will be transformed as follows:
//	fun foo(a: str)
//	fun bar(a: u32, b: u32) => u32
//
// Compiler considers such usage of tuples redundant an may issue a warning for
// such cases. Code formatting tool should transform these to canonical form.
//
// Tuple usage is restricted in the language. For example it is not possible to
// define a custom type from a tuple:
//
//	type Bar => (u32, bool)
//
// Such definition will be rejected by the compiler. Use structs if you need a
// type with multiple fields.
//
// Formal definition:
//
//	Tuple => "(" { TypeSpec "," } ")"
type Tuple struct {
	// Could be nil or contain only one element when the tree is parsed from
	// source code.
	Types []TypeSpec

	Pin source.Pin
}

var _ TypeSpec = Tuple{}

func (Tuple) Kind() tsk.Kind {
	return tsk.Tuple
}

func (t Tuple) Span() source.Span {
	return source.Span{Pin: t.Pin}
}

func (t Tuple) String() string {
	var g Printer
	g.Tuple(t)
	return g.Output()
}
