package genc

import "github.com/mebyus/ku/internal/ku/stg"

func (g *genbuf) funsig(name string, export bool, f *stg.FunDef) {
	if !export {
		g.puts("static ")
	}
	if f.Result == nil {
		g.puts("void")
	} else {
		g.typespec(f.Result)
	}
	g.line()

	g.puts(name)

	if len(f.Inputs) == 0 && f.Receiver == nil {
		g.puts("(void)")
		return
	}

	g.puts("(")
	var start int
	if f.Receiver != nil {
		g.typespec(f.Receiver)
		g.space()
		g.putb('g')
	} else {
		g.typespec(f.Inputs[0])
		g.space()
		g.puts(g.getName(f.Params[0]))
		start = 1
	}
	for i := start; i < len(f.Inputs); i += 1 {
		g.puts(", ")
		g.typespec(f.Inputs[i])
		g.space()
		g.puts(g.getName(f.Params[i]))
	}
	g.puts(")")
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
