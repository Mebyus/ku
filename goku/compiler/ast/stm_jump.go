package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/stk"
	"github.com/mebyus/ku/goku/compiler/sm"
)

// JumpNext represents jump to next loop iteration (continue) statement.
//
// Formal definition:
//
//	JumpNext => "jump" "@.next";
type JumpNext struct {
	Pin sm.Pin
}

var _ Statement = JumpNext{}

func (JumpNext) Kind() stk.Kind {
	return stk.JumpNext
}

func (j JumpNext) Span() sm.Span {
	return sm.Span{Pin: j.Pin}
}

func (j JumpNext) String() string {
	var g Printer
	g.JumpNext(j)
	return g.Output()
}

// JumpOut represents jump out of the loop (break) statement.
//
// Formal definition:
//
//	JumpOut => "jump" "@.out";
type JumpOut struct {
	Pin sm.Pin
}

var _ Statement = JumpOut{}

func (JumpOut) Kind() stk.Kind {
	return stk.JumpOut
}

func (j JumpOut) Span() sm.Span {
	return sm.Span{Pin: j.Pin}
}

func (j JumpOut) String() string {
	var g Printer
	g.JumpOut(j)
	return g.Output()
}
