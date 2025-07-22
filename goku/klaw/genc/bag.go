package genc

import "github.com/mebyus/ku/goku/compiler/ast"

func (g *Gen) RegBag(r ast.RegBag) {
	g.puts("static const ")
	g.puts(r.BagName.Str)
	g.puts("Tab ")

	g.puts(r.BagName.Str)
	g.puts("_")
	g.puts(r.Name.Str)
	g.puts("_tab = {")
	g.nl()
	g.inc()

	g.indent()
	g.puts(".type_id = ")
	id := g.State.GetTypeId(r.Name.Str)
	g.putn(id)
	g.puts(",")
	g.nl()

	for _, f := range r.Tab {
		g.indent()

		g.puts(".")
		g.puts(f.Name.Str)
		g.puts(" = (")
		g.puts(r.BagName.Str)
		g.puts("Fun")
		g.puts(f.Name.Str)
		g.puts(")((void*)(")
		g.puts(f.Fun.Str)
		g.puts(")),")

		g.nl()
	}

	g.dec()
	g.puts("};")
}
