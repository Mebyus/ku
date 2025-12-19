package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/tsk"
	"github.com/mebyus/ku/goku/compiler/sm"
)

// AnyType represents a type specifier which denotes any type.
//
// Formal definition:
//
//	AnyType -> "type"
type AnyType struct {
	Pin sm.Pin
}

var _ TypeSpec = AnyType{}

func (AnyType) Kind() tsk.Kind {
	return tsk.Type
}

func (t AnyType) Span() sm.Span {
	return sm.Span{Pin: t.Pin, Len: 4}
}

func (t AnyType) String() string {
	return "type"
}
