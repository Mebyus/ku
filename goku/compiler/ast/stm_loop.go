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
	Exp Exp
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

// ForRange represents for loop over a range of integers.
//
// Formal definition:
//
//	ForRange => "for" Name ":" TypeSpec "in" "range" "(" Exp ")" Block
//	Name => word
type ForRange struct {
	Body Block
	Name Word
	Exp  Exp
	Type TypeSpec
}

var _ Statement = ForRange{}

func (ForRange) Kind() stk.Kind {
	return stk.ForRange
}

func (r ForRange) Span() srcmap.Span {
	return r.Name.Span()
}

func (r ForRange) String() string {
	var g Printer
	g.ForRange(r)
	return g.Output()
}

// ForRange represents for loop over a range of integers.
//
// Formal definition:
//
//	ForRange2 => "for" Name ":" TypeSpec "in" "range" "(" Exp "," Exp ")" Block
//	Name => word
type ForRange2 struct {
	Body  Block
	Name  Word
	Start Exp
	End   Exp
	Type  TypeSpec
}

var _ Statement = ForRange2{}

func (ForRange2) Kind() stk.Kind {
	return stk.ForRange2
}

func (r ForRange2) Span() srcmap.Span {
	return r.Name.Span()
}

func (r ForRange2) String() string {
	var g Printer
	g.ForRange2(r)
	return g.Output()
}
