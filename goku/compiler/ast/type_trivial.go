package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/tsk"
	"github.com/mebyus/ku/goku/compiler/srcmap"
)

// Void represents "void" usage as zero-size type specifier.
// Should only be used for defining custom types.
type Void struct {
	Pin srcmap.Pin
}

var _ TypeSpec = Void{}

func (Void) Kind() tsk.Kind {
	return tsk.Void
}

func (t Void) Span() srcmap.Span {
	return srcmap.Span{Pin: t.Pin, Len: 2}
}

func (Void) String() string {
	return "void"
}
