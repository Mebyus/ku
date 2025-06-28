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
	// case DotName:
	// 	g.DotName(e)
	// case Dirty:
	// 	g.Dirty(e)
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
	// case Pack:
	// 	g.Pack(e)
	// case Chain:
	// 	g.Chain(e)
	// case Call:
	// 	g.Call(e)
	// case Ref:
	// 	g.Ref(e)
	// case Slice:
	// 	g.Slice(e)
	// case Tweak:
	// 	g.Tweak(e)
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

func (g *Gen) Symbol(n ast.Symbol) {
	g.puts(n.Name)
}

func (g *Gen) Index(x ast.Index) {
	g.puts("[")
	g.Exp(x.Exp)
	g.puts("]")
}

func (g *Gen) DerefIndex(x ast.DerefIndex) {
	g.puts(".[")
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
