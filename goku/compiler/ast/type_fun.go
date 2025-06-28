package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/tsk"
	"github.com/mebyus/ku/goku/compiler/srcmap"
)

// FunType represents function type specifier.
//
// Formal definition:
//
//	FunType => "fun" Signature
type FunType struct {
	Signature
}

var _ TypeSpec = FunType{}

func (FunType) Kind() tsk.Kind {
	return tsk.Fun
}

func (f FunType) Span() srcmap.Span {
	return srcmap.Span{Pin: 0} // TODO: place something here
}

func (f FunType) String() string {
	var g Printer
	g.FunType(f)
	return g.Output()
}
