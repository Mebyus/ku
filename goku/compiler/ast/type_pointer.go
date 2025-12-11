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

// Ref represents reference type specifier.
//
// Formal definition:
//
// Ref => "&" TypeSpec
type Ref struct {
	// Type to which pointer refers to.
	Type TypeSpec
}

var _ TypeSpec = Ref{}

func (Ref) Kind() tsk.Kind {
	return tsk.Ref
}

func (p Ref) Span() srcmap.Span {
	return p.Type.Span()
}

func (p Ref) String() string {
	var g Printer
	g.Ref(p)
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

var _ TypeSpec = AnyPointer{}

func (AnyPointer) Kind() tsk.Kind {
	return tsk.AnyPointer
}

func (p AnyPointer) Span() srcmap.Span {
	return srcmap.Span{Pin: p.Pin}
}

func (p AnyPointer) String() string {
	return "*any"
}

// AnyRef represents pointer to an unknown type.
//
// Formal definition:
//
//	AnyRef => "&" "any"
type AnyRef struct {
	Pin srcmap.Pin
}

var _ TypeSpec = AnyRef{}

func (AnyRef) Kind() tsk.Kind {
	return tsk.AnyRef
}

func (p AnyRef) Span() srcmap.Span {
	return srcmap.Span{Pin: p.Pin}
}

func (p AnyRef) String() string {
	return "&any"
}

// ArrayPointer represents array pointer type specifier.
//
// Formal definition:
//
//	ArrayPointer => "[*]" TypeSpec
type ArrayPointer struct {
	// Type to which pointer refers to.
	Type TypeSpec
}

var _ TypeSpec = ArrayPointer{}

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

// ArrayRef represents array reference type specifier.
//
// Formal definition:
//
//	ArrayRef => "[&]" TypeSpec
type ArrayRef struct {
	// Type to which pointer refers to.
	Type TypeSpec
}

var _ TypeSpec = ArrayRef{}

func (ArrayRef) Kind() tsk.Kind {
	return tsk.ArrayRef
}

func (p ArrayRef) Span() srcmap.Span {
	return p.Type.Span()
}

func (p ArrayRef) String() string {
	var g Printer
	g.ArrayRef(p)
	return g.Output()
}
