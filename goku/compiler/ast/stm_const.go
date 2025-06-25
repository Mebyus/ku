package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/stk"
	"github.com/mebyus/ku/goku/compiler/enums/tnk"
	"github.com/mebyus/ku/goku/compiler/source"
)

type TopConst struct {
	Const

	Traits
}

var _ Top = TopConst{}

func (TopConst) Kind() tnk.Kind {
	return tnk.Const
}

// Const represents constant definition statement.
//
// Formal definition:
//
//	Const => "const" Name ":" TypeSpec "=" Exp ";"
type Const struct {
	Name Word

	// Can be nil if constant type is not specified.
	Type TypeSpec

	// Specifies constant init value expression.
	//
	// Always not nil.
	Exp Exp
}

var _ Statement = Const{}

func (Const) Kind() stk.Kind {
	return stk.Const
}

func (c Const) Span() source.Span {
	return c.Name.Span()
}

func (c Const) String() string {
	var g Printer
	g.Const(c)
	return g.Output()
}
