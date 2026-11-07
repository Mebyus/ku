package ast

import (
	"github.com/mebyus/ku/internal/fbuf"
	"github.com/mebyus/ku/internal/ku/sx"
)

// Block represents a sequence of statements inside a block.
// Block can be a statement, function body, branch body, etc.
//
// Formal definition:
//
//	Block -> "{" { Statement } "}"
type Block struct {
	stm

	// List of all statements inside the block in order as they appear
	// in source code.
	//
	// Equals nil if block contains no statements.
	Nodes []Statement

	// Opening brace pin of this block.
	Pin sx.Pin
}

// Explicit interface implementation check.
var _ Statement = &Block{}

func (b *Block) Format(buf *fbuf.Buffer) {
	if len(b.Nodes) == 0 {
		buf.Puts("{}")
		return
	}

	buf.Putb('{')
	buf.Line()
	buf.Inc()
	for _, s := range b.Nodes {
		buf.Indent()
		s.Format(buf)
		buf.Line()
	}
	buf.Dec()
	buf.Indent()
	buf.Putb('}')
}

// Statement node that represents statement of any kind.
type Statement interface {
	// Discriminator method for interface implementations.
	// Only serves as a trick to enhance Go typechecking in
	// type assertions.
	//
	// Does nothing when called.
	_stm()

	// Format statement as text that should be very close
	// or identical to what it was parsed from.
	//
	// Used mainly for tests and debugging.
	Format(*fbuf.Buffer)
}

// Embed this to quickly implement _stm() discriminator from Statement interface.
// Do not use it for anything else.
type stm struct{}

func (stm) _stm() {}

// ZoneBlock represents a code block with specified zone name.
type ZoneBlock struct {
	stm

	// List of all statements inside the block in order as they appear
	// in source code.
	//
	// Equals nil if block contains no statements.
	Nodes []Statement

	// Zone name.
	Name string

	// Zone name pin.
	Pin sx.Pin
}

func (z *ZoneBlock) Format(buf *fbuf.Buffer) {
	buf.Puts("zone ")
	buf.Puts(z.Name)
	buf.Space()

	b := Block{Nodes: z.Nodes}
	b.Format(buf)
}
