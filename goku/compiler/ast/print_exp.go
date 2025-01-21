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
	case Nil:
		g.Nil(e)
	case Unary:
		g.Unary(e)
	case Binary:
		g.Binary(e)
	case Paren:
		g.Paren(e)
	case Pack:
		g.Pack(e)
	case Chain:
		g.Chain(e)
	case Call:
		g.Call(e)
	case Ref:
		g.Ref(e)
	case Slice:
		g.Slice(e)
	default:
		panic(fmt.Sprintf("unexpected \"%s\" (=%d) expression (%T)", e.Kind(), e.Kind(), e))
	}
}

func (g *Printer) Slice(s Slice) {
	g.Chain(s.Chain)
	g.puts("[")
	if s.Start != nil {
		g.Exp(s.Start)
	}
	g.puts(":")
	if s.End != nil {
		g.Exp(s.End)
	}
	g.puts("]")
}

func (g *Printer) Ref(r Ref) {
	g.Chain(r.Chain)
	g.puts(".&")
}

func (g *Printer) Call(c Call) {
	g.Chain(c.Chain)
	g.puts("(")
	g.Args(c.Args)
	g.puts(")")
}

func (g *Printer) Args(args []Exp) {
	if len(args) == 0 {
		return
	}

	g.Exp(args[0])
	for _, arg := range args[1:] {
		g.puts(", ")
		g.Exp(arg)
	}
}

func (g *Printer) Chain(c Chain) {
	g.puts(c.Start.Str)

	for _, p := range c.Parts {
		g.Part(p)
	}
}

func (g *Printer) Part(p Part) {
	switch p := p.(type) {
	case Index:
		g.Index(p)
	case Select:
		g.Select(p)
	case Deref:
		g.Deref(p)
	case SelectTest:
		g.SelectTest(p)
	default:
		panic(fmt.Sprintf("unexpected \"%s\" (=%d) chain part (%T)", p.Kind(), p.Kind(), p))
	}
}

func (g *Printer) SelectTest(s SelectTest) {
	g.puts(".test.")
	g.puts(s.Name.Str)
}

func (g *Printer) Deref(d Deref) {
	g.puts(".@")
}

func (g *Printer) Select(s Select) {
	g.puts(".")
	g.puts(s.Name.Str)
}

func (g *Printer) Index(x Index) {
	g.puts("[")
	g.Exp(x.Exp)
	g.puts("]")
}

func (g *Printer) DerefIndex(x DerefIndex) {
	g.puts(".[")
	g.Exp(x.Exp)
	g.puts("]")
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

func (g *Printer) Nil(n Nil) {
	g.puts("nil")
}

func (g *Printer) Dirty(d Dirty) {
	g.puts("dirty")
}
