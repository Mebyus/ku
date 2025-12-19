package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/tsk"
	"github.com/mebyus/ku/goku/compiler/sm"
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

func (p Pointer) Span() sm.Span {
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
// Ref -> "&" TypeSpec
type Ref struct {
	// Type to which pointer refers to.
	Type TypeSpec
}

var _ TypeSpec = Ref{}

func (Ref) Kind() tsk.Kind {
	return tsk.Ref
}

func (p Ref) Span() sm.Span {
	return p.Type.Span()
}

func (p Ref) String() string {
	var g Printer
	g.Ref(p)
	return g.Output()
}

// VoidPointer represents pointer to an unknown type.
//
// Formal definition:
//
//	VoidPointer -> "*" "void"
type VoidPointer struct {
	Pin sm.Pin
}

var _ TypeSpec = VoidPointer{}

func (VoidPointer) Kind() tsk.Kind {
	return tsk.VoidPointer
}

func (p VoidPointer) Span() sm.Span {
	return sm.Span{Pin: p.Pin}
}

func (p VoidPointer) String() string {
	return "*void"
}

// VoidRef represents pointer to an unknown type.
//
// Formal definition:
//
//	VoidRef => "&" "void"
type VoidRef struct {
	Pin sm.Pin
}

var _ TypeSpec = VoidRef{}

func (VoidRef) Kind() tsk.Kind {
	return tsk.VoidRef
}

func (p VoidRef) Span() sm.Span {
	return sm.Span{Pin: p.Pin}
}

func (p VoidRef) String() string {
	return "&void"
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

func (p ArrayPointer) Span() sm.Span {
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
//	ArrayRef -> "[&]" TypeSpec
type ArrayRef struct {
	// Type to which pointer refers to.
	Type TypeSpec
}

var _ TypeSpec = ArrayRef{}

func (ArrayRef) Kind() tsk.Kind {
	return tsk.ArrayRef
}

func (p ArrayRef) Span() sm.Span {
	return p.Type.Span()
}

func (p ArrayRef) String() string {
	var g Printer
	g.ArrayRef(p)
	return g.Output()
}
