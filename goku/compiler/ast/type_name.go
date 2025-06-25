package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/tsk"
	"github.com/mebyus/ku/goku/compiler/srcmap"
)

// Formal definition:
//
//	TypeName => word
type TypeName struct {
	Name Word
}

// Explicit interface implementation check.
var _ TypeSpec = TypeName{}

func (TypeName) Kind() tsk.Kind {
	return tsk.Name
}

func (n TypeName) Span() srcmap.Span {
	return n.Name.Span()
}

func (n TypeName) String() string {
	return n.Name.Str
}
