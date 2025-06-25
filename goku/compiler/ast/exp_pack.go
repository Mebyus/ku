package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/exk"
	"github.com/mebyus/ku/goku/compiler/srcmap"
)

// Pack represents an expression formed by a list of expressions. Although pack
// is classified as expression, it cannot appear in all places where expression
// is expected. A complete list of examples where pack usage is allowed:
//
//	// Multiple assignment from a call
//	a, err := foo();
//
//	// Multiple assignment from a pack
//	a, b, c = c, a, b;
//
//	// Return statement inside a function with multiple return values
//	ret a + 1, b;
//
// Parser forms a pack only when at least two expressions are present in a list.
//
// Note that list of call arguments is a separate language construct. The example
// below does not have a pack expression:
//
//	foo(1, a, "hello")
//
// Examples of contexts where pack cannot be used:
//
//	// Index expression
//	a.foo[i]
//	a.foo[1, i] // invalid
//
//	// Any subexpression
//	a + (1 - b)
//	a + (b, 1 - b) // invalid
//
// Formal definition:
//
//	Pack => Exp "," Exp { "," Exp } // trailing comma is not allowed
type Pack struct {
	nodeExp

	// Always contains at least two elements.
	List []Exp
}

var _ Exp = Pack{}

func (Pack) Kind() exk.Kind {
	return exk.Pack
}

func (p Pack) Span() srcmap.Span {
	return srcmap.Span{Pin: p.List[0].Span().Pin}
}

func (p Pack) String() string {
	var g Printer
	g.Pack(p)
	return g.Output()
}
