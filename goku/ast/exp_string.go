package ast

import (
	"github.com/mebyus/ku/goku/enums/exk"
	"github.com/mebyus/ku/goku/source"
)

type String struct {
	// String literal value represented by token.
	Val string

	Pin source.Pin
}

var _ Exp = String{}

func (String) Kind() exk.Kind {
	return exk.String
}

func (s String) Span() source.Span {
	return source.Span{Pin: s.Pin, Len: uint32(len(s.Val)) + 2}
}

func (s String) String() string {
	var g Printer
	g.String(s)
	return g.Str()
}
