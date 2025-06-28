package genc

import (
	"fmt"

	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/enums/tnk"
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
		panic("not implemented")
	case tnk.Var:
		g.TopVar(text.Variables[i])
	case tnk.Type:
		g.Type(text.Types[i])
	case tnk.Test:
		panic("not implemented")
	case tnk.Method:
		panic("not supported")
	case tnk.FunStub:
		panic("not implemented")
	case tnk.Gen:
		panic("not supported")
	case tnk.GenBind:
		panic("not supported")
	default:
		panic(fmt.Sprintf("unexpected \"%s\" (=%d) top level node (i=%d)", k, k, i))
	}
	g.nl()
}
