package genc

import "github.com/mebyus/ku/goku/compiler/ast"

func (g *Gen) Fun(f ast.Fun) {
	g.FunHead(f.Name.Str, f.Signature)
	g.space()
	g.Block(f.Body)
}

func (g *Gen) FunStub(s ast.FunStub) {
	g.FunHead(s.Name.Str, s.Signature)
	g.semi()
}

func getTestFunName(name string) string {
	return "run_test_" + name
}

func (g *Gen) TestFun(t ast.Fun) {
	if !g.test() {
		return
	}

	t.Name.Str = getTestFunName(t.Name.Str)
	g.Fun(t)
}

func (g *Gen) FunHead(name string, s ast.Signature) {
	g.puts("static ")
	if s.Never {
		g.puts("_Noreturn ")
	}
	if s.Result == nil {
		g.puts("void")
	} else {
		g.TypeSpec(s.Result)
	}
	g.nl()

	g.puts(name)
	g.FunParams(s.Params)
}

func (g *Gen) FunParams(params []ast.Param) {
	if len(params) == 0 {
		g.puts("(void)")
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
