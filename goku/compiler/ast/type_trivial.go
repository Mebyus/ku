package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/tsk"
	"github.com/mebyus/ku/goku/compiler/srcmap"
)

type Trivial struct {
	Pin srcmap.Pin
}

var _ TypeSpec = Trivial{}

func (Trivial) Kind() tsk.Kind {
	return tsk.Trivial
}

func (t Trivial) Span() srcmap.Span {
	return srcmap.Span{Pin: t.Pin, Len: 2}
}

func (Trivial) String() string {
	return "()"
}
