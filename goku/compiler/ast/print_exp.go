package ast

import (
	"fmt"

	"github.com/mebyus/ku/goku/compiler/char"
)

func (g *Printer) Exp(exp Exp) {
	switch e := exp.(type) {
	case Name:
		g.Name(e)
	case Dirty:
		g.Dirty(e)
	case Integer:
		g.Integer(e)
	case String:
		g.String(e)
	case Binary:
		g.Binary(e)
	case Paren:
		g.Paren(e)
	default:
		panic(fmt.Sprintf("unexpected \"%s\" (=%d) expression", e.Kind(), e.Kind()))
	}
}

func (g *Printer) Binary(b Binary) {
	g.Exp(b.A)
	g.space()
	g.puts(b.Op.Kind.String())
	g.space()
	g.Exp(b.B)
}

func (g *Printer) Paren(p Paren) {
	g.puts("(")
	g.Exp(p.Exp)
	g.puts(")")
}

func (g *Printer) Name(n Name) {
	g.puts(n.Word)
}

func (g *Printer) Integer(n Integer) {
	g.puts(n.String())
}

func (g *Printer) String(s String) {
	g.puts("\"")
	g.puts(char.Escape(s.Val))
	g.puts("\"")
}

func (g *Printer) Dirty(d Dirty) {
	g.puts("dirty")
}
