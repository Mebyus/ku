package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/tsk"
	"github.com/mebyus/ku/goku/compiler/source"
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

func (n TypeName) Span() source.Span {
	return n.Name.Span()
}

func (n TypeName) String() string {
	return n.Name.Str
}
