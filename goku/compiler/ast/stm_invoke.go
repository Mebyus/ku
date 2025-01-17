package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/stk"
	"github.com/mebyus/ku/goku/compiler/source"
)

// Invoke represents a call expression statement.
//
// Formal definition:
//
//	Invoke => Call ";"
type Invoke struct {
	Call Call
}

var _ Statement = Invoke{}

func (Invoke) Kind() stk.Kind {
	return stk.Invoke
}

func (i Invoke) Span() source.Span {
	return i.Call.Span()
}

func (i Invoke) String() string {
	var g Printer
	g.Invoke(i)
	return g.Output()
}
