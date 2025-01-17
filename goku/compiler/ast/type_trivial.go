package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/tsk"
	"github.com/mebyus/ku/goku/compiler/source"
)

type Trivial struct {
	Pin source.Pin
}

var _ TypeSpec = Trivial{}

func (Trivial) Kind() tsk.Kind {
	return tsk.Trivial
}

func (t Trivial) Span() source.Span {
	return source.Span{Pin: t.Pin, Len: 2}
}

func (Trivial) String() string {
	return "()"
}
