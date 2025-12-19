package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/stk"
	"github.com/mebyus/ku/goku/compiler/sm"
)

// If represents if-else statement.
//
// Formal definition:
//
//	If => IfClause { ElseIfClause } [ ElseClause ]
//
//	IfClause     => "if" Exp Block
//	ElseIfClause => "else" IfClause
//	ElseClause   => "else" Block
type If struct {
	If IfClause

	// Equals nil if there are no "else if" clauses in statement
	ElseIfs []IfClause

	// Equals nil if there is no "else" clause in statement
	Else *Block
}

type IfClause struct {
	// Branch condition. Always not nil.
	Exp Exp

	Body Block
}

var _ Statement = If{}

func (If) Kind() stk.Kind {
	return stk.If
}

func (i If) Span() sm.Span {
	return i.If.Exp.Span()
}

func (i If) String() string {
	var g Printer
	return g.Output()
}

type StaticIf struct {
	If IfClause

	// Equals nil if there are no "else if" clauses in statement
	ElseIfs []IfClause

	// Equals nil if there is no "else" clause in statement
	Else *Block
}

var _ Statement = StaticIf{}

func (StaticIf) Kind() stk.Kind {
	return stk.StaticIf
}

func (i StaticIf) Span() sm.Span {
	return i.If.Exp.Span()
}

func (i StaticIf) String() string {
	var g Printer
	return g.Output()
}
