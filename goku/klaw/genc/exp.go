package genc

import (
	"fmt"

	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/char"
)

func (g *Gen) Exp(exp ast.Exp) {
	switch e := exp.(type) {
	case nil:
		panic("nil exp")
	case ast.Symbol:
		g.Symbol(e)
	case ast.DotName:
		panic("not supported")
	case ast.Dirty:
		panic("not supported")
	case ast.Integer:
		g.Integer(e)
	case ast.String:
		g.String(e)
	case ast.Nil:
		g.Nil(e)
	case ast.Unary:
		g.Unary(e)
	case ast.Binary:
		g.Binary(e)
	case ast.Paren:
		g.Paren(e)
	case ast.Pack:
		panic("not supported")
	case ast.Chain:
		g.Chain(e)
	case ast.Call:
		g.Call(e)
	case ast.Ref:
		g.Ref(e)
	case ast.Slice:
		panic("not supported")
	case ast.Tweak:
		panic("not supported")
	default:
		panic(fmt.Sprintf("unexpected \"%s\" (=%d) expression (%T)", e.Kind(), e.Kind(), e))
	}
}

func (g *Gen) Binary(b ast.Binary) {
	g.Exp(b.A)
	g.space()
	g.puts(b.Op.Kind.String())
	g.space()
	g.Exp(b.B)
}

func (g *Gen) Unary(u ast.Unary) {
	g.puts(u.Op.Kind.String())
	g.Exp(u.Exp)
}

func (g *Gen) Paren(p ast.Paren) {
	g.puts("(")
	g.Exp(p.Exp)
	g.puts(")")
}

func (g *Gen) Ref(r ast.Ref) {
	g.puts("&")
	g.Chain(r.Chain)
}

func (g *Gen) Call(c ast.Call) {
	g.Chain(c.Chain)
	g.puts("(")
	g.Args(c.Args)
	g.puts(")")
}

func (g *Gen) Args(args []ast.Exp) {
	if len(args) == 0 {
		return
	}

	g.Exp(args[0])
	for _, arg := range args[1:] {
		g.puts(", ")
		g.Exp(arg)
	}
}

func (g *Gen) Chain(c ast.Chain) {
	g.chain(c.Start, c.Parts)
}

func (g *Gen) chain(start ast.Word, parts []ast.Part) {
	if len(parts) == 0 {
		g.puts(start.Str)
		return
	}

	i := len(parts) - 1
	last := parts[i]
	rest := parts[:i]

	switch p := last.(type) {
	case ast.Index:
		g.chain(start, rest)
		g.Index(p)
	case ast.Deref:
		g.putb('*')
		g.chain(start, rest)
	case ast.DerefIndex:
		g.chain(start, rest)
		g.DerefIndex(p)
	case ast.Select:
		g.chain(start, rest)
		g.Select(p)
	default:
		panic(fmt.Sprintf("unexpected \"%s\" (=%d) chain part (%T)", p.Kind(), p.Kind(), p))
	}
}

func (g *Gen) Symbol(n ast.Symbol) {
	g.puts(n.Name)
}

func (g *Gen) Select(s ast.Select) {
	g.puts(".")
	g.puts(s.Name.Str)
}

func (g *Gen) Index(x ast.Index) {
	g.puts("[")
	g.Exp(x.Exp)
	g.puts("]")
}

func (g *Gen) DerefIndex(x ast.DerefIndex) {
	g.puts("[")
	g.Exp(x.Exp)
	g.puts("]")
}

func (g *Gen) Integer(n ast.Integer) {
	g.puts(n.String())
}

func (g *Gen) String(s ast.String) {
	if len(s.Val) == 0 {
		g.puts("empty_str")
		return
	}
	g.puts("make_str((u8*)(u8\"")
	g.puts(char.Escape(s.Val))
	g.puts("\"), ")
	g.putlen(len(s.Val))
	g.putb(')')
}

func (g *Gen) Nil(n ast.Nil) {
	g.puts("nil")
}
