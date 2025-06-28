package genc

import "github.com/mebyus/ku/goku/compiler/ast"

func (g *Gen) Statement(s ast.Statement) {

}

func (g *Gen) Block(b ast.Block) {
	if len(b.Nodes) == 0 {
		g.puts("{}")
		return
	}

	g.puts("{")
	g.nl()
	g.inc()

	for _, n := range b.Nodes {
		g.indent()
		g.Statement(n)
		g.nl()
	}

	g.dec()
	g.indent()
	g.puts("}")
}


func (g *Gen) TopVar(v ast.TopVar) {
	g.puts("static ")
	g.TypeSpec(v.Type)
	g.nl()
	g.puts(v.Name.Str)

	if v.Exp == nil {
		g.semi()
		return
	}

	g.puts(" = ")
	g.Exp(v.Exp)
	g.semi()
}
