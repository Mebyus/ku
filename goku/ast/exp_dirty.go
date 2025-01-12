package ast

import (
	"github.com/mebyus/ku/goku/enums/exk"
	"github.com/mebyus/ku/goku/source"
)

// Dirty represents usage of "dirty" keyword as expression.
type Dirty struct {
	Pin source.Pin
}

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
