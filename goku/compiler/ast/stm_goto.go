package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/stk"
	"github.com/mebyus/ku/goku/compiler/sm"
)

type Goto struct {
	// Label name.
	Name string

	Pin sm.Pin
}

var _ Statement = Goto{}

func (Goto) Kind() stk.Kind {
	return stk.Goto
}

func (t Goto) Span() sm.Span {
	return sm.Span{Pin: t.Pin}
}

func (t Goto) String() string {
	var g Printer
	g.Goto(t)
	return g.Output()
}

// Gonext represents jump to next loop iteration (continue) statement.
//
// Formal definition:
//
//	Gonext => "gonext" ";"
type Gonext struct {
	Pin sm.Pin
}

var _ Statement = Gonext{}

func (Gonext) Kind() stk.Kind {
	return stk.Gonext
}

func (n Gonext) Span() sm.Span {
	return sm.Span{Pin: n.Pin}
}

func (n Gonext) String() string {
	var g Printer
	g.Gonext(n)
	return g.Output()
}

// Break represents jump out of the loop (break) statement.
//
// Formal definition:
//
//	Break => "break" ";"
type Break struct {
	Pin sm.Pin
}

var _ Statement = Break{}

func (Break) Kind() stk.Kind {
	return stk.Break
}

func (b Break) Span() sm.Span {
	return sm.Span{Pin: b.Pin}
}

func (b Break) String() string {
	var g Printer
	g.Break(b)
	return g.Output()
}
