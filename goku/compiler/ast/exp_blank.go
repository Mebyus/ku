package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/exk"
	"github.com/mebyus/ku/goku/compiler/sm"
)

// Blank represents usage of "_" as expression.
//
// When used as assignment target tells the compiler that corresponding
// value is skipped. For example:
//
//	n, _ = parse_int("42"); // error is skipped
//	_ = p.next(); // we need only function call, not its result
type Blank struct {
	nodeExp

	Pin sm.Pin
}

// Explicit interface implementation check.
var _ Exp = Blank{}

func (Blank) Kind() exk.Kind {
	return exk.Blank
}

func (b Blank) Span() sm.Span {
	return sm.Span{Pin: b.Pin, Len: 1}
}

func (Blank) String() string {
	return "_"
}
