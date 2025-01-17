package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/tsk"
	"github.com/mebyus/ku/goku/compiler/source"
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

func (p Pointer) Span() source.Span {
	return p.Type.Span()
}

func (p Pointer) String() string {
	var g Printer
	g.Pointer(p)
	return g.Output()
}
