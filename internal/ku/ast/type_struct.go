package ast

import (
	"github.com/mebyus/ku/internal/fbuf"
	"github.com/mebyus/ku/internal/ku/sx"
)

// Struct represents struct type specifier.
//
// Formal definition:
//
//	Struct -> "struct" "{" { Field [ "," ] } "}"
type Struct struct {
	spec

	// Can be nil (if struct does not have fields).
	Fields []Field

	Pin sx.Pin
}

var _ TypeSpec = &Struct{}

func (s *Struct) Format(buf *fbuf.Buffer) {
	if len(s.Fields) == 0 {
		buf.Puts("struct {}")
		return
	}

	buf.Puts("struct {\n")
	buf.Inc()
	for _, f := range s.Fields {
		buf.Indent()
		buf.Puts(f.Name)
		buf.Puts(": ")
		f.Type.Format(buf)
		buf.Line()
	}
	buf.Dec()
	buf.Indent()
	buf.Putb('}')
}

// Field represents a single field in struct or union, or form type specifier.
//
// Formal definition:
//
//	Field  -> Name ":" TypeSpec
//	Name   -> word
type Field struct {
	Name string
	Type TypeSpec

	// pin of field name
	Pin sx.Pin
}
