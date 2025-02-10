package ast

import (
	"fmt"

	"github.com/mebyus/ku/goku/compiler/enums/tnk"
)

func (g *Printer) Gen(gen Gen) {
	g.puts("gen ")
	g.puts(gen.Name.Str)

	g.puts("(")
	g.Params(gen.Params)
	g.puts(")")

	if gen.Control != nil {
		g.Static(*gen.Control)
	}
}

func (g *Printer) GenBind(b GenBind) {
	g.puts("gen ")
	g.puts(b.Name.Str)

	g.puts("(...) ")
	g.GenBlock(&b.Body)
}

func (g *Printer) GenBlock(b *GenBlock) {
	if len(b.OrderIndex) == 0 {
		g.puts("{}")
		return
	}

	g.puts("{")
	g.nl()
	g.nl()
	g.genByIndex(b, b.OrderIndex[0])
	for _, x := range b.OrderIndex[1:] {
		g.nl()
		g.genByIndex(b, x)
	}
	g.nl()
	g.puts("}")
}

func (g *Printer) genByIndex(b *GenBlock, x NodeIndex) {
	k := x.Kind
	i := x.Index
	switch k {
	case 0:
		panic(fmt.Sprintf("unspecified top level node (i=%d)", i))
	case tnk.Fun:
		g.Fun(b.Functions[i])
	case tnk.Const:
		g.TopConst(b.Constants[i])
	case tnk.Type:
		g.Type(b.Types[i])
	case tnk.Method:
		g.Method(b.Methods[i])
	default:
		panic(fmt.Sprintf("unexpected \"%s\" (=%d) top level node (i=%d)", k, k, i))
	}
	g.nl()
}
