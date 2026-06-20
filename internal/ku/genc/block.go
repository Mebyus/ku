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
	case *stg.Branch:
		g.branch(s)
	case *stg.LineIf:
		g.lineif(s)
	case *stg.While:
		g.while(s)
	case *stg.Block:
		g.block(s)
	default:
		panic(fmt.Sprintf("unexpected %T statement", s))
	}
}

func (g *Buffer) lineif(f *stg.LineIf) {
	g.puts("if (")
	g.exp(f.Exp)
	g.puts(") ")
	g.node(f.Then)
}

func (g *Buffer) branch(f *stg.Branch) {
	g.puts("if (")
	g.exp(f.Exp)
	g.puts(") ")
	g.block(&f.Body)
	if f.Else == nil {
		return
	}

	g.puts(" else ")
	g.block(f.Else)
}

func (g *Buffer) while(w *stg.While) {
	g.puts("while (")
	g.exp(w.Exp)
	g.puts(") ")
	g.block(&w.Body)
}

func (g *Buffer) ret(r *stg.Return) {
	if r.Exp == nil {
		g.puts("return;")
		return
	}

	g.puts("return ")
	g.exp(r.Exp)
	g.semi()
}
