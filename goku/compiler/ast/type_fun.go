package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/tsk"
	"github.com/mebyus/ku/goku/compiler/sm"
)

// FunType represents function type specifier.
//
// Formal definition:
//
//	FunType -> "fun" Signature
type FunType struct {
	Signature

	Pin sm.Pin
}

var _ TypeSpec = FunType{}

func (FunType) Kind() tsk.Kind {
	return tsk.Fun
}

func (f FunType) Span() sm.Span {
	return sm.Span{Pin: f.Pin}
}

func (f FunType) String() string {
	var g Printer
	g.FunType(f)
	return g.Output()
}
