package ast

import (
	"fmt"

	"github.com/mebyus/ku/goku/compiler/char"
)

func (g *Printer) Exp(exp Exp) {
	switch e := exp.(type) {
	case Symbol:
		g.Symbol(e)
	case Dirty:
		g.Dirty(e)
	case Integer:
		g.Integer(e)
	case String:
		g.String(e)
	case Unary:
		g.Unary(e)
	case Binary:
		g.Binary(e)
	case Paren:
		g.Paren(e)
	case Pack:
		g.Pack(e)
	default:
		panic(fmt.Sprintf("unexpected \"%s\" (=%d) expression (%T)", e.Kind(), e.Kind(), e))
	}
}

func (g *Printer) Pack(p Pack) {
	g.Exp(p.List[0])
	for _, exp := range p.List[1:] {
		g.puts(", ")
		g.Exp(exp)
	}
}

func (g *Printer) Binary(b Binary) {
	g.Exp(b.A)
	g.space()
	g.puts(b.Op.Kind.String())
	g.space()
	g.Exp(b.B)
}

func (g *Printer) Unary(u Unary) {
	g.puts(u.Op.Kind.String())
	g.Exp(u.Exp)
}

func (g *Printer) Paren(p Paren) {
	g.puts("(")
	g.Exp(p.Exp)
	g.puts(")")
}

func (g *Printer) Symbol(n Symbol) {
	g.puts(n.Name)
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
