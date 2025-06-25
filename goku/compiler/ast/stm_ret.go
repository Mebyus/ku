package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/stk"
	"github.com/mebyus/ku/goku/compiler/srcmap"
)

// Formal definition:
//
//	Ret => "ret" [ Exp ] ";"
type Ret struct {
	// Equals nil if return does not have expression.
	Exp Exp

	Pin srcmap.Pin
}

var _ Statement = Ret{}

func (Ret) Kind() stk.Kind {
	return stk.Ret
}

func (r Ret) Span() srcmap.Span {
	return srcmap.Span{Pin: r.Pin}
}

func (r Ret) String() string {
	var g Printer
	g.Ret(r)
	return g.Output()
}
