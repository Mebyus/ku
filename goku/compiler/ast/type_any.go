package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/tsk"
	"github.com/mebyus/ku/goku/compiler/source"
)

// AnyType represents a type specifier which denotes any type.
//
// Formal definition:
//
//	AnyType => "type"
type AnyType struct {
	Pin source.Pin
}

var _ TypeSpec = AnyType{}

func (AnyType) Kind() tsk.Kind {
	return tsk.Type
}

func (t AnyType) Span() source.Span {
	return source.Span{Pin: t.Pin, Len: 4}
}

func (t AnyType) String() string {
	return "type"
}
