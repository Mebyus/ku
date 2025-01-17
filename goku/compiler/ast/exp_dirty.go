package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/exk"
	"github.com/mebyus/ku/goku/compiler/source"
)

// Dirty represents usage of "dirty" keyword as expression.
type Dirty struct {
	nodeExp

	Pin source.Pin
}

// Explicit interface implementation check.
var _ Exp = Dirty{}

func (Dirty) Kind() exk.Kind {
	return exk.Dirty
}

func (d Dirty) Span() source.Span {
	return source.Span{Pin: d.Pin, Len: uint32(len(d.String()))}
}

func (Dirty) String() string {
	return "dirty"
}
