package ast

import (
	"github.com/mebyus/ku/internal/fbuf"
	"github.com/mebyus/ku/internal/ku/sx"
)

// Loop represents a loop with no condition.
//
// Formal definition:
//
//	Loop -> "for" Block
type Loop struct {
	stm

	Body Block

	Pin sx.Pin
}

var _ Statement = &Loop{}

func (p *Loop) Format(buf *fbuf.Buffer) {
	buf.Puts("for ")
	p.Body.Format(buf)
}

// While represents a loop with condition.
//
// Formal definition:
//
//	While -> "for" Exp Block
type While struct {
	stm

	Body Block

	// Loop condition. Always not nil.
	Exp Exp

	Pin sx.Pin
}

var _ Statement = &While{}

func (w *While) Format(buf *fbuf.Buffer) {
	buf.Puts("for ")
	w.Exp.Format(buf)
	buf.Space()
	w.Body.Format(buf)
}

// ForNum represents a for loop that iterates specified
// number of times.
type ForNum struct {
	stm

	Body Block

	// Range start. Can be nil if omitted.
	A Exp

	// Range end. Always not nil.
	B Exp

	Pin sx.Pin
}

var _ Statement = &ForNum{}

func (f *ForNum) Format(buf *fbuf.Buffer) {
	buf.Puts("for [")
	if f.A != nil {
		f.A.Format(buf)
	}
	buf.Putb(':')
	f.B.Format(buf)
	buf.Puts("] ")
	f.Body.Format(buf)
}

// ForRange represents a for loop with variable that
// iterates over specified range.
type ForRange struct {
	stm

	Body Block

	// Range start. Can be nil if omitted.
	A Exp

	// Range end. Always not nil.
	B Exp

	// Loop variable name.
	Name string

	// Name pin.
	Pin sx.Pin
}

var _ Statement = &ForRange{}

func (f *ForRange) Format(buf *fbuf.Buffer) {
	buf.Puts("for ")
	buf.Puts(f.Name)
	buf.Puts(" := [")
	if f.A != nil {
		f.A.Format(buf)
	}
	buf.Putb(':')
	f.B.Format(buf)
	buf.Puts("] ")
	f.Body.Format(buf)
}
