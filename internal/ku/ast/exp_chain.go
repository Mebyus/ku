package ast

import (
	"github.com/mebyus/ku/internal/fbuf"
	"github.com/mebyus/ku/internal/ku/sx"
)

// Chain represents an operand which starts with identifier and has at least
// one chain part operator attached to it (select, call, deref, index, etc.)
type Chain struct {
	operand

	// name of identifier which starts the chain
	Name string

	// Always has at least one element.
	Parts []Part

	// pin of start identifier
	Pin sx.Pin
}

var _ Operand = &Chain{}

func (c *Chain) Format(buf *fbuf.Buffer) {
	buf.Puts(c.Name)
	for _, p := range c.Parts {
		p.Format(buf)
	}
}

type Part interface {
	// Discriminator method for interface implementations.
	// Only serves as a trick to enhance Go typechecking in
	// type assertions.
	//
	// Does nothing when called.
	_part()

	// Format chain part as text that should be very close
	// or identical to what it was parsed from.
	//
	// Used mainly for tests and debugging.
	Format(*fbuf.Buffer)
}

// Embed this to quickly implement _part() discriminator from Part interface.
// Do not use it for anything else.
type part struct{}

func (part) _part() {}

// Select named field from chain operand.
type Select struct {
	part

	Name string

	Pin sx.Pin
}

// Explicit interface implementation check.
var _ Part = &Select{}

func (s *Select) Format(buf *fbuf.Buffer) {
	buf.Putb('.')
	buf.Puts(s.Name)
}

// Index into chain operand.
type Index struct {
	part

	Exp Exp

	Pin sx.Pin
}

var _ Part = &Index{}

func (x *Index) Format(buf *fbuf.Buffer) {
	buf.Putb('[')
	x.Exp.Format(buf)
	buf.Putb(']')
}
