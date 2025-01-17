package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/tsk"
	"github.com/mebyus/ku/goku/compiler/source"
)

// Struct represents struct type specifier.
//
// Formal definition:
//
//	Struct => "struct" "{" { Field "," } "}"
type Struct struct {
	// Can be nil (if struct does not have fields).
	Fields []Field

	Pin source.Pin
}

var _ TypeSpec = Struct{}

func (Struct) Kind() tsk.Kind {
	return tsk.Struct
}

func (s Struct) Span() source.Span {
	return source.Span{Pin: s.Pin}
}

func (s Struct) String() string {
	var g Printer
	g.Struct(s)
	return g.Output()
}

// Field represents a single field in struct type specifier.
//
// Formal definition:
//
//	Field  => Name ":" TypeSpec
//	Name   => word
type Field struct {
	Name Word
	Type TypeSpec
}
