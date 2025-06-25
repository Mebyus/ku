package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/stk"
	"github.com/mebyus/ku/goku/compiler/enums/tnk"
	"github.com/mebyus/ku/goku/compiler/srcmap"
)

type TopAlias struct {
	Alias

	Traits
}

var _ Top = TopAlias{}

func (TopAlias) Kind() tnk.Kind {
	return tnk.Alias
}

// Alias represents alias definition statement.
//
// Formal definition:
//
//	Alias => "let" Name "=>" Exp ";"
type Alias struct {
	Name Word

	// Specifies aliased expression.
	//
	// Always not nil.
	Exp Exp
}

var _ Statement = Alias{}

func (Alias) Kind() stk.Kind {
	return stk.Alias
}

func (a Alias) Span() srcmap.Span {
	return a.Name.Span()
}

func (a Alias) String() string {
	var g Printer
	g.Alias(a)
	return g.Output()
}
