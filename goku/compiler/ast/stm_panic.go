package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/stk"
	"github.com/mebyus/ku/goku/compiler/srcmap"
)

// Formal definition:
//
//	Panic => "panic" "(" String ")" ";"
type Panic struct {
	Msg string

	Pin srcmap.Pin
}

var _ Statement = Panic{}

func (Panic) Kind() stk.Kind {
	return stk.Panic
}

func (p Panic) Span() srcmap.Span {
	return srcmap.Span{Pin: p.Pin}
}

func (p Panic) String() string {
	var g Printer
	g.Panic(p)
	return g.Output()
}
