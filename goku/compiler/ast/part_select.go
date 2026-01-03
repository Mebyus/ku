package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/exk"
	"github.com/mebyus/ku/goku/compiler/sm"
)

type Select struct {
	nodePart

	Name Word
}

// Explicit interface implementation check.
var _ Part = Select{}

func (Select) Kind() exk.Kind {
	return exk.Select
}

func (s Select) Span() sm.Span {
	return s.Name.Span()
}

func (s Select) String() string {
	var g Printer
	g.Select(s)
	return g.Output()
}

// BagSelect represents chain part.
//
// Formal definition:
//
//	BagSelect -> ".(" Name ")"
//
//	Name -> word
type BagSelect struct {
	nodePart

	Name Word
}

// Explicit interface implementation check.
var _ Part = BagSelect{}

func (BagSelect) Kind() exk.Kind {
	return exk.BagSelect
}

func (s BagSelect) Span() sm.Span {
	return s.Name.Span()
}

func (s BagSelect) String() string {
	var g Printer
	g.BagSelect(s)
	return g.Output()
}

// SelectTest represents chain part of selecting a test function from import.
//
// Formal definition:
//
//	SelectTest => "." "test" "." Name
//
//	Name => word
type SelectTest struct {
	nodePart

	Name Word
}

// Explicit interface implementation check.
var _ Part = SelectTest{}

func (SelectTest) Kind() exk.Kind {
	return exk.SelectTest
}

func (s SelectTest) Span() sm.Span {
	return s.Name.Span()
}

func (s SelectTest) String() string {
	var g Printer
	g.SelectTest(s)
	return g.Output()
}
