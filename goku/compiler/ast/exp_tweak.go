package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/exk"
	"github.com/mebyus/ku/goku/compiler/sm"
)

// Tweak represents usage of tweak as expression.
//
// Formal definition:
//
//	Tweak => Chain ".{" { ObjField "," } "}"
type Tweak struct {
	nodeOperand

	Chain Chain

	Fields []ObjField
}

// Explicit interface implementation check.
var _ Operand = Tweak{}

func (Tweak) Kind() exk.Kind {
	return exk.Tweak
}

func (t Tweak) Span() sm.Span {
	return t.Chain.Span()
}

func (t Tweak) String() string {
	var g Printer
	g.Tweak(t)
	return g.Output()
}

type ObjField struct {
	Name Word
	Exp  Exp
}
