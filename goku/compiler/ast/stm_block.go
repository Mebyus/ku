package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/stk"
	"github.com/mebyus/ku/goku/compiler/source"
)

type Block struct {
	Nodes []Statement

	// Opening brace pin of this block.
	Pin source.Pin
}

var _ Statement = Block{}

func (Block) Kind() stk.Kind {
	return stk.Block
}

func (b Block) Span() source.Span {
	return source.Span{Pin: b.Pin}
}

func (b Block) String() string {
	var g Printer
	g.Block(b)
	return g.Output()
}
