package genc

import "github.com/mebyus/ku/goku/compiler/ast"

func (g *Gen) Fun(f ast.Fun) {
	g.puts("static ")
	if f.Signature.Never {
		g.puts("_Noreturn ")
	}
	if f.Signature.Result == nil {
		g.puts("void")
	} else {
		g.TypeSpec(f.Signature.Result)
	}
	g.nl()

	g.puts(f.Name.Str)
	g.FunParams(f.Signature.Params)
	g.space()
	g.Block(f.Body)
}

func (g *Gen) FunParams(params []ast.Param) {
	if len(params) == 0 {
		g.puts("()")
		return
	}

	g.puts("(")
	g.Param(params[0])
	for _, p := range params[1:] {
		g.puts(", ")
		g.Param(p)
	}
	g.puts(")")
}

func (g *Gen) Param(p ast.Param) {
	g.TypeSpec(p.Type)
	g.space()
	g.puts(p.Name.Str)
}
