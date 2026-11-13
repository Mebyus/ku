package genc

import "github.com/mebyus/ku/internal/ku/stg"

func (g *genbuf) funsig(name string, export bool, f *stg.FunDef) {
	if !export {
		g.puts("static ")
	}
	if f.Result == nil {
		g.puts("void")
	} else {
		g.typ(f.Result)
	}
	g.line()

	g.puts(name)

	if len(f.Inputs) == 0 {
		g.puts("(void)")
	} else {
		g.puts("(")
		g.typ(f.Inputs[0])
		g.space()
		g.puts(g.getName(f.Params[0]))
		for i := 1; i < len(f.Inputs); i += 1 {
			g.puts(", ")
			g.typ(f.Inputs[i])
			g.space()
			g.puts(g.getName(f.Params[i]))
		}
		g.puts(")")
	}
}

func (g *genbuf) fun(s *stg.Symbol) {
	f := s.Def.(*stg.FunDef)
	g.funsig(g.getName(s), s.IsExport(), f)
	g.space()
	g.block(&f.Body)
	g.line()
}

func (g *genbuf) funstub(s *stg.Symbol) {
	g.funsig(g.getName(s), s.IsExport(), s.Def.(*stg.FunDef))
	g.semi()
	g.line()
}
