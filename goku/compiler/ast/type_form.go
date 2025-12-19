package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/tsk"
	"github.com/mebyus/ku/goku/compiler/sm"
)

// Form represents form type specifier.
//
// Formal definition:
//
//	Form => "(" { Field "," } ")"
type Form struct {
	// Always not nil when produced by parser.
	Fields []Field
}

var _ TypeSpec = Form{}

func (Form) Kind() tsk.Kind {
	return tsk.Form
}

func (f Form) Span() sm.Span {
	return sm.Span{Pin: f.Fields[0].Name.Pin}
}

func (f Form) String() string {
	var g Printer
	g.Form(f)
	return g.Output()
}
