package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/tsk"
	"github.com/mebyus/ku/goku/compiler/srcmap"
)

// Pointer represents pointer type specifier.
//
// Formal definition:
//
//	Pointer => "*" TypeSpec
type Pointer struct {
	// Type to which pointer refers to.
	Type TypeSpec
}

var _ TypeSpec = Pointer{}

func (Pointer) Kind() tsk.Kind {
	return tsk.Pointer
}

func (p Pointer) Span() srcmap.Span {
	return p.Type.Span()
}

func (p Pointer) String() string {
	var g Printer
	g.Pointer(p)
	return g.Output()
}

// AnyPointer represents pointer to an unknown type.
//
// Formal definition:
//
//	AnyPointer => "*" "any"
type AnyPointer struct {
	Pin srcmap.Pin
}

var _ TypeSpec = Pointer{}

func (AnyPointer) Kind() tsk.Kind {
	return tsk.AnyPointer
}

func (p AnyPointer) Span() srcmap.Span {
	return srcmap.Span{Pin: p.Pin}
}

func (p AnyPointer) String() string {
	return "*any"
}

// ArrayPointer represents array pointer type specifier.
//
// Formal definition:
//
//	Pointer => "[*]" TypeSpec
type ArrayPointer struct {
	// Type to which pointer refers to.
	Type TypeSpec
}

var _ TypeSpec = Pointer{}

func (ArrayPointer) Kind() tsk.Kind {
	return tsk.ArrayPointer
}

func (p ArrayPointer) Span() srcmap.Span {
	return p.Type.Span()
}

func (p ArrayPointer) String() string {
	var g Printer
	g.ArrayPointer(p)
	return g.Output()
}
