package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/stk"
	"github.com/mebyus/ku/goku/compiler/srcmap"
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

func (l Loop) Span() srcmap.Span {
	return l.Body.Span()
}

func (l Loop) String() string {
	var g Printer
	g.Loop(l)
	return g.Output()
}

// While represents a loop with condition.
//
// Formal definition:
//
//	While => "for" Exp Block
type While struct {
	Body Block

	// Loop condition. Always not nil.
	Exp  Exp
}


var _ Statement = While{}

func (While) Kind() stk.Kind {
	return stk.Loop
}

func (w While) Span() srcmap.Span {
	return w.Body.Span()
}

func (w While) String() string {
	var g Printer
	g.While(w)
	return g.Output()
}
