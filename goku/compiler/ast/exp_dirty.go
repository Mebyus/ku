package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/exk"
	"github.com/mebyus/ku/goku/compiler/source"
)

// Dirty represents usage of "?" as expression.
//
// When used as init expression tells the compiler to skip variable
// initialization (even with default value). When used as return expression
// tells the compiler that returned value will not be used by the caller and
// thus anything can be placed in this value at runtime.
//
// Formal definition:
//
//	Dirty => "?"
type Dirty struct {
	nodeExp

	Pin source.Pin
}

// Explicit interface implementation check.
var _ Exp = Dirty{}

func (Dirty) Kind() exk.Kind {
	return exk.Dirty
}

func (d Dirty) Span() source.Span {
	return source.Span{Pin: d.Pin, Len: uint32(len(d.String()))}
}

func (Dirty) String() string {
	return "?"
}
