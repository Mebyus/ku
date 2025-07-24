package genc

import (
	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/char"
)

func (g *Gen) typedefBagFun(bag string, f ast.BagFun) {
	g.puts("typedef ")

	s := f.Signature
	if s.Result == nil {
		g.puts("void")
	} else {
		g.TypeSpec(s.Result)
	}

	g.puts(" (*")
	g.puts(bag)
	g.puts("Fun_")
	g.puts(f.Name.Str)
	g.puts(")")

	g.puts("(uint")
	for _, p := range f.Signature.Params {
		g.puts(", ")
		g.TypeSpec(p.Type)
	}
	g.puts(");")
}

func (g *Gen) typedefBagType(name string, b ast.Bag) {
	if len(b.Funs) == 0 {
		panic("empty bag type")
	}

	for _, f := range b.Funs {
		g.typedefBagFun(name, f)
	}

	g.nl()
	g.puts("typedef struct {")
	g.nl()
	g.inc()

	g.indent()
	g.puts("uint type_id;")
	g.nl()

	for _, f := range b.Funs {
		g.indent()
		g.puts(name)
		g.puts("Fun_")
		g.puts(f.Name.Str)
		g.space()
		g.puts(f.Name.Str)
		g.semi()
		g.nl()
	}

	g.dec()
	g.puts("} ")
	g.puts(name)
	g.puts("Tab;")

	g.nl()
	g.puts("typedef struct {")
	g.nl()
	g.inc()

	g.indent()
	g.puts("uint obj;")
	g.nl()

	g.indent()
	g.puts("const ")
	g.puts(name)
	g.puts("Tab* tab;")
	g.nl()

	g.dec()
	g.puts("} ")
	g.puts(name)
	g.semi()
}

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
		g.puts("Fun_")
		g.puts(f.Name.Str)
		g.puts(")((void*)(")
		g.puts(f.Fun.Str)
		g.puts(")),")

		g.nl()
	}

	g.dec()
	g.puts("};")
	g.nl()

	g.nl()
	g.puts("static ")
	g.puts(r.BagName.Str)
	g.nl()
	g.puts("make_")
	g.puts(char.SnakeCase(r.BagName.Str))
	g.puts("_from_")
	g.puts(r.Name.Str)

	g.puts("(")
	g.TypeSpec(r.Type)
	g.space()
	g.puts(r.Name.Str)
	g.puts(") {")
	g.nl()
	g.inc()

	g.indent()
	g.puts("static_assert(sizeof(")
	g.TypeSpec(r.Type)
	g.puts(") <= sizeof(uint));")
	g.nl()

	g.indent()
	g.puts(r.BagName.Str)
	g.puts(" b = {")
	g.nl()
	g.inc()

	g.indent()
	g.puts(".obj = (uint)(")
	g.puts(r.Name.Str)
	g.puts("),")
	g.nl()

	g.indent()
	g.puts(".tab = &")
	g.puts(r.BagName.Str)
	g.puts("_")
	g.puts(r.Name.Str)
	g.puts("_tab,")
	g.nl()
	g.dec()
	g.indent()
	g.puts("};")
	g.nl()

	g.indent()
	g.puts("return b;")
	g.nl()

	g.dec()
	g.puts("}")
}
