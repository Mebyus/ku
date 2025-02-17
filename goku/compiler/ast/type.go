package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/tnk"
	"github.com/mebyus/ku/goku/compiler/source"
)

// Type represents top level type definition construct.
//
// Formal definitino:
//
//	Type => [ "pub" ] "type" Name [ "in" BagList ] "=>" TypeSpec
//	BagList => "(" { BagName "," } ")"
//	Name => word
//	BagName => word
type Type struct {
	Name Word

	// Optional list of bags which this type must fit into.
	Bags []Word

	Spec TypeSpec

	Traits
}

var _ Top = Type{}

func (Type) Kind() tnk.Kind {
	return tnk.Type
}

func (t Type) Span() source.Span {
	return t.Name.Span()
}

func (t Type) String() string {
	var g Printer
	g.Type(t)
	return g.Output()
}
