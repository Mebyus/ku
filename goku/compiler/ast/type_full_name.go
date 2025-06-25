package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/tsk"
	"github.com/mebyus/ku/goku/compiler/srcmap"
)

// Formal definition:
//
//	TypeFullName => ImportName "." TypeName
//	ImportName   => word
type TypeFullName struct {
	Import Word
	Name   Word
}

// Explicit interface implementation check
var _ TypeSpec = TypeFullName{}

func (TypeFullName) Kind() tsk.Kind {
	return tsk.FullName
}

func (n TypeFullName) Span() srcmap.Span {
	return srcmap.Span{
		Pin: n.Import.Pin,
		Len: uint32(len(n.Import.Str)+len(n.Name.Str)) + 1,
	}
}

func (n TypeFullName) String() string {
	var g Printer
	g.TypeFullName(n)
	return g.Output()
}
