package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/stk"
	"github.com/mebyus/ku/goku/compiler/srcmap"
)

// Formal definition:
//
//	StaticMust => "#must" "(" Exp ")" ";"
type StaticMust struct {
	// Always not nil.
	Exp Exp
}

var _ Statement = StaticMust{}

func (StaticMust) Kind() stk.Kind {
	return stk.StaticMust
}

func (m StaticMust) Span() srcmap.Span {
	return m.Exp.Span()
}

func (m StaticMust) String() string {
	var g Printer
	g.StaticMust(m)
	return g.Output()
}
