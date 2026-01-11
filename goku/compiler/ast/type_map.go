package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/tsk"
	"github.com/mebyus/ku/goku/compiler/sm"
)

// Formal definition:
//
//	Map -> "map" "(" Key "," Value ")"
//	Key -> TypeSpec
//	Value -> TypeSpec
type Map struct {
	Key   TypeSpec
	Value TypeSpec
}

// Explicit interface implementation check.
var _ TypeSpec = Map{}

func (Map) Kind() tsk.Kind {
	return tsk.Map
}

func (m Map) Span() sm.Span {
	return m.Key.Span()
}

func (m Map) String() string {
	panic("not implemented")
}
