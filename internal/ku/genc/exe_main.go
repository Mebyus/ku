package genc

import "github.com/mebyus/ku/internal/ku/stg"

func (g *genbuf) genMain(s *stg.Symbol) {
	def := s.Def.(*stg.FunDef)
	if len(def.Params) != 0 || def.Result != nil {
		panic("TODO")
	}

	g.line()
	g.puts("int main() {\n")
	g.inc()

	g.indent()
	g.puts(g.getName(s))
	g.puts("();\n")

	g.dec()
	g.puts("}\n")
}
