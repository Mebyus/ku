package ast

import (
	"fmt"
	"io"

	"github.com/mebyus/ku/goku/compiler/enums/tnk"
)

func Render(w io.Writer, text *Text) error {
	var p Printer
	p.Text(text)
	_, err := p.WriteTo(w)
	return err
}

func (g *Printer) Text(text *Text) {
	if len(text.OrderIndex) == 0 {
		return
	}

	g.topByIndex(text, text.OrderIndex[0])
	for _, x := range text.OrderIndex[1:] {
		g.nl()
		g.topByIndex(text, x)
	}
}

func (g *Printer) topByIndex(text *Text, x TopNodeIndex) {
	k := x.Kind
	i := x.Index
	switch k {
	case 0:
		panic(fmt.Sprintf("unspecified top level node (i=%d)", i))
	case tnk.Fun:
		g.Fun(text.Functions[i])
	case tnk.Let:
		g.TopLet(text.Constants[i])
	case tnk.Var:
		g.TopVar(text.Variables[i])
	case tnk.Type:
		g.Type(text.Types[i])
	case tnk.Test:
		g.Test(text.Tests[i])
	case tnk.Method:
		g.Method(text.Methods[i])
	case tnk.FunStub:
		g.FunStub(text.FunStubs[i])
	default:
		panic(fmt.Sprintf("unexpected \"%s\" (=%d) top level node (i=%d)", k, k, i))
	}
	g.nl()
}
