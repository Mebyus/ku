package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/stk"
	"github.com/mebyus/ku/goku/compiler/enums/tnk"
	"github.com/mebyus/ku/goku/compiler/sm"
)

// TopVar represents top level variable definition.
//
// Formal definition:
//
//	Var => "var" Name ":" TypeSpec [ "=" Exp ] ";"
type TopVar struct {
	Var

	Traits
}

var _ Top = TopVar{}

func (TopVar) Kind() tnk.Kind {
	return tnk.Var
}

// Var represents variable definition statement.
//
// Formal definition:
//
//	Var => "var" Name ":" TypeSpec [ "=" Exp ] ";"
type Var struct {
	Name Word

	Type TypeSpec

	// Specifies variable init value expression.
	//
	// Equals nil if init expression is empty.
	// In that case default init value is used when variable is created.
	Exp Exp
}

var _ Statement = Var{}

func (Var) Kind() stk.Kind {
	return stk.Var
}

func (v Var) Span() sm.Span {
	return v.Name.Span()
}

func (v Var) String() string {
	var g Printer
	g.Var(v)
	return g.Output()
}
