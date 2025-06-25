package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/stk"
	"github.com/mebyus/ku/goku/compiler/srcmap"
)

// Lookup represents lookup statement.
// Creates a "wormhole" into another scope. Code after this statement will use
// the "wormhole" to lookup symbols inside connected scope before falling back
// to default lookup mechanism.
//
// Formal definition:
//
//	Lookup => "#lookup" Exp ";"
type Lookup struct {
	// Expression that specifies the scope.
	Exp Exp

	Pin srcmap.Pin
}

var _ Statement = Lookup{}

func (Lookup) Kind() stk.Kind {
	return stk.Lookup
}

func (l Lookup) Span() srcmap.Span {
	return srcmap.Span{Pin: l.Pin}
}

func (l Lookup) String() string {
	var g Printer
	g.Lookup(l)
	return g.Output()
}
