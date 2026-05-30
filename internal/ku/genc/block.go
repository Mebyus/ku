package genc

import (
	"fmt"

	"github.com/mebyus/ku/internal/ku/stg"
)

func (g *Buffer) block(block *stg.Block) {
	if len(block.Nodes) == 0 {
		g.puts("{}")
		return
	}

	g.puts("{")
	g.nl()
	g.inc()
	for _, s := range block.Nodes {
		g.indent()
		g.node(s)
		g.nl()
	}
	g.dec()
	g.indent()
	g.puts("}")
}

func (g *Buffer) node(s stg.Statement) {
	switch s := s.(type) {
	case *stg.Return:
		g.ret(s)
	default:
		panic(fmt.Sprintf("unexpected %T statement", s))
	}
}

func (g *Buffer) ret(r *stg.Return) {
	if r.Exp == nil {
		g.puts("return;")
		return
	}

	g.puts("return ")
	panic("stub")
	g.semi()
}
