package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/stk"
	"github.com/mebyus/ku/goku/compiler/sm"
)

// Formal definition:
//
//	Panic => "panic" "(" String ")" ";"
type Panic struct {
	Msg string

	Pin sm.Pin
}

var _ Statement = Panic{}

func (Panic) Kind() stk.Kind {
	return stk.Panic
}

func (p Panic) Span() sm.Span {
	return sm.Span{Pin: p.Pin}
}

func (p Panic) String() string {
	var g Printer
	g.Panic(p)
	return g.Output()
}
