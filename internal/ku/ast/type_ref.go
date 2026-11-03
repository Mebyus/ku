package ast

import (
	"github.com/mebyus/ku/internal/fbuf"
	"github.com/mebyus/ku/internal/ku/sx"
)

// Ref represents reference type specifier.
//
// Formal definition:
//
//	Ref -> "&" TypeSpec
type Ref struct {
	spec

	// Type to which pointer refers to.
	Type TypeSpec

	Pin sx.Pin
}

var _ TypeSpec = &Ref{}

func (r *Ref) Format(buf *fbuf.Buffer) {
	buf.Putb('&')
	r.Type.Format(buf)
}
