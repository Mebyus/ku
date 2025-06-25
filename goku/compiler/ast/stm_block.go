package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/stk"
	"github.com/mebyus/ku/goku/compiler/srcmap"
)

// Block represents a sequence of statements inside a block.
// Block can be a statement, function body, branch body, etc.
//
// Formal definition:
//
//	Block => "{" { Statement } "}"
type Block struct {
	Nodes []Statement

	// Opening brace pin of this block.
	Pin srcmap.Pin
}

var _ Statement = Block{}

func (Block) Kind() stk.Kind {
	return stk.Block
}

func (b Block) Span() srcmap.Span {
	return srcmap.Span{Pin: b.Pin}
}

func (b Block) String() string {
	var g Printer
	g.Block(b)
	return g.Output()
}

// Debug represents a block statement under debug compile-time condition.
//
// Formal definition:
//
//	Debug => "#debug" Block
type Debug struct {
	Block Block
}

var _ Statement = Debug{}

func (Debug) Kind() stk.Kind {
	return stk.Debug
}

func (d Debug) Span() srcmap.Span {
	return d.Block.Span()
}

func (d Debug) String() string {
	var g Printer
	g.Debug(d)
	return g.Output()
}

// Static represents a block statement executed at compile-time.
//
// Formal definition:
//
//	Static => "#{" { Statement } "}"
type Static struct {
	Nodes []Statement

	// Opening brace pin of this block.
	Pin srcmap.Pin
}
