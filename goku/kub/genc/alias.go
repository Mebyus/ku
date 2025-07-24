package genc

import "github.com/mebyus/ku/goku/compiler/ast"

func (g *Gen) Alias(a ast.TopAlias) {
	g.puts("#define ")
	g.puts(a.Name.Str)
	g.puts(" (")
	g.Exp(a.Exp)
	g.puts(")")
}
