package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/stk"
	"github.com/mebyus/ku/goku/compiler/source"
)

// Stub represents stub statement.
// It means that something in source code is not implemented yet.
// Program execution will panic on this statement.
//
// Formal definition:
//
//	Stub => "#stub" ";"
type Stub struct {
	Pin source.Pin
}

var _ Statement = Stub{}

func (Stub) Kind() stk.Kind {
	return stk.Stub
}

func (s Stub) Span() source.Span {
	return source.Span{Pin: s.Pin}
}

func (s Stub) String() string {
	var g Printer
	g.Stub(s)
	return g.Output()
}
