package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/exk"
	"github.com/mebyus/ku/goku/compiler/source"
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

func (s Select) Span() source.Span {
	return s.Name.Span()
}

func (s Select) String() string {
	var g Printer
	g.Select(s)
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

func (s SelectTest) Span() source.Span {
	return s.Name.Span()
}

func (s SelectTest) String() string {
	var g Printer
	g.SelectTest(s)
	return g.Output()
}
