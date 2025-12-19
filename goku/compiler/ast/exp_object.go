package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/exk"
	"github.com/mebyus/ku/goku/compiler/sm"
)

// Object represents usage of object literal as expression.
//
// Formal definition:
//
//	Object => "{" { ObjField "," } "}"
type Object struct {
	nodeOperand

	Fields []ObjField

	Pin sm.Pin
}

// Explicit interface implementation check.
var _ Operand = Object{}

func (Object) Kind() exk.Kind {
	return exk.Object
}

func (o Object) Span() sm.Span {
	return sm.Span{Pin: o.Pin}
}

func (o Object) String() string {
	var g Printer
	g.Object(o)
	return g.Output()
}
