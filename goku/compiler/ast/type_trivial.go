package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/tsk"
	"github.com/mebyus/ku/goku/compiler/sm"
)

// Void represents "void" usage as zero-size type specifier.
// Should only be used for defining custom types.
type Void struct {
	Pin sm.Pin
}

var _ TypeSpec = Void{}

func (Void) Kind() tsk.Kind {
	return tsk.Void
}

func (t Void) Span() sm.Span {
	return sm.Span{Pin: t.Pin, Len: 2}
}

func (Void) String() string {
	return "void"
}
