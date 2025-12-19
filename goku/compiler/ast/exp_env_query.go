package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/exk"
	"github.com/mebyus/ku/goku/compiler/sm"
)

// EnvQuery represents compile-time query of env variable or constant.
//
// FormalDefinition:
//
//	EnvQuery => "#:" { word "." }
type EnvQuery struct {
	nodeOperand

	Name string
	Pin  sm.Pin
}

// Explicit interface implementation check.
var _ Operand = EnvQuery{}

func (EnvQuery) Kind() exk.Kind {
	return exk.EnvQuery
}

func (q EnvQuery) Span() sm.Span {
	return sm.Span{Pin: q.Pin, Len: uint32(len(q.Name))}
}

func (q EnvQuery) String() string {
	var g Printer
	g.EnvQuery(q)
	return g.Output()
}
