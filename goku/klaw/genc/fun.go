package genc

import (
	"fmt"

	"github.com/mebyus/ku/goku/compiler/ast"
)

func (g *Gen) Fun(f ast.Fun) {
	linkType, ok := getPropValue(f.Traits, "link.type")
	if ok && linkType == "c.main" {
		g.cmainFunHead()
	} else {
		g.FunHead(f.Name.Str, f.Signature)
	}
	g.space()
	g.Block(f.Body)
}

func (g *Gen) FunStub(s ast.FunStub) {
	linkType, ok := getPropValue(s.Traits, "link.type")
	if ok {
		if linkType == "external" {
			g.externalFunStub(s)
			return
		}
		panic(fmt.Sprintf("unexpected value \"%s\" of link.type property", linkType))
	}

	g.FunHead(s.Name.Str, s.Signature)
	g.semi()
}

func (g *Gen) cmainFunHead() {
	g.puts("int main(int argc, char** argv, char** envp)")
}

func (g *Gen) externalFunStub(s ast.FunStub) {
	linkName, ok := getPropValue(s.Traits, "link.name")
	if !ok {
		g.funDeclaration(s.Name.Str, s.Signature)
		g.semi()
		return
	}
	if linkName == s.Name.Str {
		panic(fmt.Sprintf("property link.name has the same value \"%s\" as function stub name", linkName))
	}

	g.funDeclaration(linkName, s.Signature)
	g.semi()
	g.nl()
	g.nl()

	g.FunHead(s.Name.Str, s.Signature)
	g.puts(" {")
	g.nl()
	g.inc()
	g.level += 1

	g.indent()
	if s.Signature.Result != nil {
		g.puts("return ")
	}
	g.puts(linkName)
	g.puts("(")
	g.paramNamesList(s.Signature.Params)
	g.puts(");")
	g.nl()

	g.level -= 1
	g.dec()
	g.puts("}")
}

func (g *Gen) paramNamesList(params []ast.Param) {
	if len(params) == 0 {
		return
	}

	g.puts(params[0].Name.Str)
	for _, p := range params[1:] {
		g.puts(", ")
		g.puts(p.Name.Str)
	}
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
	g.funDeclaration(name, s)
}

func (g *Gen) funDeclaration(name string, s ast.Signature) {
	g.funResult(s)
	g.nl()

	g.puts(name)
	g.FunParams(s.Params)
}

func (g *Gen) funResult(s ast.Signature) {
	if s.Never {
		g.puts("_Noreturn ")
	}
	if s.Result == nil {
		g.puts("void")
	} else {
		g.TypeSpec(s.Result)
	}
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
