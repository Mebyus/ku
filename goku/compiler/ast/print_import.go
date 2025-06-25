package ast

import "github.com/mebyus/ku/goku/compiler/srcmap/origin"

func (g *Printer) ImportBlocks(blocks []ImportBlock) {
	if len(blocks) == 0 {
		return
	}

	if !g.empty() {
		g.nl()
	}
	g.ImportBlock(blocks[0])
	g.nl()
	for _, block := range blocks[1:] {
		g.nl()
		g.ImportBlock(block)
		g.nl()
	}
}

func (g *Printer) ImportBlock(block ImportBlock) {
	g.puts("import ")
	if block.Origin != origin.Loc {
		g.puts(block.Origin.String())
		g.space()
	}

	if len(block.Imports) == 0 {
		g.puts("{}")
		return
	}

	g.puts("{")
	g.nl()
	g.inc()

	for _, i := range block.Imports {
		g.indent()
		g.puts(i.Name.Str)
		g.puts(" => \"")
		g.puts(i.String.Str)
		g.puts("\"")
		g.nl()
	}

	g.dec()
	g.puts("}")
}
