package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/stk"
	"github.com/mebyus/ku/goku/compiler/sm"
)

type Match struct {
	Cases []MatchCase

	Exp Exp

	Else *Block
}

var _ Statement = Match{}

func (Match) Kind() stk.Kind {
	return stk.Match
}

func (m Match) Span() sm.Span {
	return m.Exp.Span()
}

func (m Match) String() string {
	// var g Printer
	// g.While(w)
	// return g.Output()
	return ""
}

type MatchCase struct {
	Body Block

	// Always has at least one element.
	List []Exp
}
