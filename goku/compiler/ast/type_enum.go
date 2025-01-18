package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/tsk"
	"github.com/mebyus/ku/goku/compiler/source"
)

type Enum struct {
	// Base integer type.
	Base TypeName

	Entries []EnumEntry
}

type EnumEntry struct {
	Name Word

	// Can be nil if entry does not have explicitly assigned value.
	Exp Exp
}

// Explicit interface implementation check.
var _ TypeSpec = Enum{}

func (Enum) Kind() tsk.Kind {
	return tsk.Enum
}

func (e Enum) Span() source.Span {
	return e.Base.Span()
}

func (e Enum) String() string {
	var g Printer
	g.Enum(e)
	return g.Output()
}
