package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/stk"
	"github.com/mebyus/ku/goku/compiler/sm"
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

func (l Loop) Span() sm.Span {
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
	Exp Exp
}

var _ Statement = While{}

func (While) Kind() stk.Kind {
	return stk.While
}

func (w While) Span() sm.Span {
	return w.Body.Span()
}

func (w While) String() string {
	var g Printer
	g.While(w)
	return g.Output()
}

// ForRange represents for loop over a range of integers.
//
// Formal definition:
//
//	ForRange => "for" Name [ ":" TypeSpec ] "=" "[" [ Exp ] ":" Exp "]" Block
//	Name => word
type ForRange struct {
	Body Block
	Name Word

	// May be nil when start is omitted. Means it equals zero value.
	Start Exp

	// Always not nil.
	End Exp

	// May be nil for auto-type loop variable.
	Type TypeSpec
}

var _ Statement = ForRange{}

func (ForRange) Kind() stk.Kind {
	return stk.ForRange
}

func (r ForRange) Span() sm.Span {
	return r.Name.Span()
}

func (r ForRange) String() string {
	var g Printer
	g.ForRange(r)
	return g.Output()
}
