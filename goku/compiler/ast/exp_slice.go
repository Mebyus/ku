package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/exk"
	"github.com/mebyus/ku/goku/compiler/sm"
)

// Formal definition:
//
//	Slice -> "[" [ Start ] ":" [ End ] "]"
type Slice struct {
	nodeOperand

	Chain Chain

	// Part before ":". Can be nil if expression is omitted.
	Start Exp

	// Part after ":". Can be nil if expression is omitted.
	End Exp
}

// Explicit interface implementation check.
var _ Operand = Slice{}

func (Slice) Kind() exk.Kind {
	return exk.Slice
}

func (s Slice) Span() sm.Span {
	if s.Start != nil {
		return s.Start.Span()
	}
	if s.End != nil {
		return s.End.Span()
	}
	return s.Chain.TailSpan()
}

func (s Slice) String() string {
	var g Printer
	g.Slice(s)
	return g.Output()
}

// Formal definition:
//
//	DerefSlice -> ".[" [ Start ] ":" [ End ] "]"
type DerefSlice struct {
	nodeOperand

	Chain Chain

	// Part before ":". Can be nil if expression is omitted.
	Start Exp

	// Part after ":". Can be nil if expression is omitted.
	End Exp
}

// Explicit interface implementation check.
var _ Operand = DerefSlice{}

func (DerefSlice) Kind() exk.Kind {
	return exk.DerefSlice
}

func (s DerefSlice) Span() sm.Span {
	if s.Start != nil {
		return s.Start.Span()
	}
	if s.End != nil {
		return s.End.Span()
	}
	return s.Chain.TailSpan()
}

func (s DerefSlice) String() string {
	var g Printer
	g.DerefSlice(s)
	return g.Output()
}
