package ast

import (
	"github.com/mebyus/ku/internal/fbuf"
	"github.com/mebyus/ku/internal/ku/sx"
)

// Span represents span type specifier.
//
// Formal definition:
//
//	Span -> "[]" TypeSpec
type Span struct {
	spec

	// Span element type specifier.
	Type TypeSpec

	Pin sx.Pin
}

// Explicit interface implementation check.
var _ TypeSpec = &Span{}

func (s *Span) Format(buf *fbuf.Buffer) {
	buf.Puts("[]")
	s.Type.Format(buf)
}
