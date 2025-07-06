package genc

import (
	"fmt"

	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/enums/tnk"
	"github.com/mebyus/ku/goku/compiler/srcmap"
)

func (g *Gen) Nodes(text *ast.Text) {
	if len(text.OrderIndex) == 0 {
		return
	}

	if !g.empty() {
		g.nl()
	}
	g.topByIndex(text, text.OrderIndex[0])
	for _, x := range text.OrderIndex[1:] {
		g.nl()
		g.topByIndex(text, x)
	}
}

func (g *Gen) topByIndex(text *ast.Text, x ast.NodeIndex) {
	k := x.Kind
	i := x.Index
	switch k {
	case 0:
		panic(fmt.Sprintf("unspecified top level node (i=%d)", i))
	case tnk.Fun:
		g.Fun(text.Functions[i])
	case tnk.Const:
		g.TopConst(text.Constants[i])
	case tnk.Var:
		g.TopVar(text.Variables[i])
	case tnk.Type:
		g.Type(text.Types[i])
	case tnk.Test:
		g.Test(text.Tests[i])
	case tnk.Method:
		panic("not supported")
	case tnk.FunStub:
		g.FunStub(text.FunStubs[i])
	case tnk.Alias:
		g.Alias(text.Aliases[i])
	case tnk.Must:
		g.StaticMust(text.Musts[i])
	case tnk.Gen:
		panic("not supported")
	case tnk.GenBind:
		panic("not supported")
	default:
		panic(fmt.Sprintf("unexpected \"%s\" (=%d) top level node (i=%d)", k, k, i))
	}
	g.nl()
}

func (g *Gen) textPosArgs(pin srcmap.Pin) {
	pos, err := g.State.Map.DecodePin(pin)
	if err != nil {
		panic(err)
	}

	g.str(pos.Path)
	g.puts(", ")
	g.putn(uint64(pos.Pos.Line) + 1)
}
