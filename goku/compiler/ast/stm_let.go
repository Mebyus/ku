package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/stk"
	"github.com/mebyus/ku/goku/compiler/enums/tnk"
	"github.com/mebyus/ku/goku/compiler/source"
)

type TopLet struct {
	Let
}

var _ Top = TopLet{}

func (TopLet) Kind() tnk.Kind {
	return tnk.Let
}

// Let represents constant definition statement.
//
// Formal definition:
//
//	Let => "let" Name ":" TypeSpec "=" Exp ";"
type Let struct {
	Name Word

	Type TypeSpec

	// Specifies constant init value expression.
	//
	// Always not nil.
	Exp Exp
}

var _ Statement = Let{}

func (Let) Kind() stk.Kind {
	return stk.Let
}

func (l Let) Span() source.Span {
	return l.Name.Span()
}

func (l Let) String() string {
	var g Printer
	g.Let(l)
	return g.Output()
}
