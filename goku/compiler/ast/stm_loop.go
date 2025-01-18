package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/stk"
	"github.com/mebyus/ku/goku/compiler/source"
)

// Loop represents a loop without condition.
//
// Formal definition:
//
//	Loop => "for" Block
type Loop struct {
	Body Block
}

var _ Statement = Loop{}

func (Loop) Kind() stk.Kind {
	return stk.Loop
}

func (l Loop) Span() source.Span {
	return l.Body.Span()
}

func (l Loop) String() string {
	var g Printer
	g.Loop(l)
	return g.Output()
}
