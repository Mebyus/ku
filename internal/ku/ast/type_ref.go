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

// ArrayPointer represents array pointer type specifier.
//
// Formal definition:
//
//	ArrayPointer -> "[*]" TypeSpec
type ArrayPointer struct {
	spec

	// Type to which array pointer refers to.
	Type TypeSpec

	Pin sx.Pin
}

var _ TypeSpec = &ArrayPointer{}

func (p *ArrayPointer) Format(buf *fbuf.Buffer) {
	buf.Puts("[*]")
	p.Type.Format(buf)
}
