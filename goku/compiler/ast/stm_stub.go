package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/stk"
	"github.com/mebyus/ku/goku/compiler/srcmap"
)

// Stub represents stub statement.
// It means that something in source code is not implemented yet.
// Program execution will panic on this statement.
//
// Formal definition:
//
//	Stub -> "stub" ";"
type Stub struct {
	Pin srcmap.Pin
}

var _ Statement = Stub{}

func (Stub) Kind() stk.Kind {
	return stk.Stub
}

func (s Stub) Span() srcmap.Span {
	return srcmap.Span{Pin: s.Pin}
}

func (s Stub) String() string {
	var g Printer
	g.Stub(s)
	return g.Output()
}

// Never represents never statement.
// It means that by programmer's logic execution must never reach this point.
// Program execution will panic on this statement.
//
// Formal definition:
//
//	Never -> "never" ";"
type Never struct {
	Pin srcmap.Pin
}

var _ Statement = Never{}

func (Never) Kind() stk.Kind {
	return stk.Never
}

func (n Never) Span() srcmap.Span {
	return srcmap.Span{Pin: n.Pin}
}

func (n Never) String() string {
	var g Printer
	g.Never(n)
	return g.Output()
}
