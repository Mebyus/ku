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
	case ast.Float:
		g.Float(e)
	case ast.String:
		g.String(e)
	case ast.Rune:
		g.Rune(e)
	case ast.True:
		g.True(e)
	case ast.False:
		g.False(e)
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
	case ast.Object:
		g.Object(e)
	case ast.List:
		g.List(e)
	case ast.Slice:
		panic("not supported")
	case ast.Tweak:
		panic("not supported")
	case ast.TypeId:
		g.TypeId(e)
	case ast.ErrorId:
		g.ErrorId(e)
	case ast.Size:
		g.Size(e)
	case ast.Cast:
		g.Cast(e)
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
	case ast.DerefSelect:
		g.chain(start, rest)
		g.DerefSelect(p)
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

func (g *Gen) DerefSelect(d ast.DerefSelect) {
	g.puts("->")
	g.puts(d.Name.Str)
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

func (g *Gen) Float(f ast.Float) {
	g.puts(f.Val)
}

func (g *Gen) String(s ast.String) {
	g.str(s.Val)
}

func (g *Gen) Rune(r ast.Rune) {
	g.puts("'")
	g.puts(char.EscapeRune(rune(r.Val)))
	g.puts("'")
}

func (g *Gen) str(s string) {
	if len(s) == 0 {
		g.puts("empty_str")
		return
	}
	g.puts("make_str((u8*)(u8\"")
	g.puts(char.Escape(s))
	g.puts("\"), ")
	g.putlen(len(s))
	g.putb(')')
}

func (g *Gen) True(t ast.True) {
	g.puts("true")
}

func (g *Gen) False(f ast.False) {
	g.puts("false")
}

func (g *Gen) Nil(n ast.Nil) {
	g.puts("nil")
}

func (g *Gen) TypeId(t ast.TypeId) {
	id := g.State.GetTypeId(t.Name.Str)
	g.putn(id)
}

func (g *Gen) ErrorId(e ast.ErrorId) {
	id := g.State.GetErrorId(e.Name.Str)
	g.putn(id)
}

func (g *Gen) Size(s ast.Size) {
	g.puts("sizeof(")
	g.TypeSpec(s.Exp)
	g.puts(")")
}

func (g *Gen) Cast(c ast.Cast) {
	g.puts("(")
	g.TypeSpec(c.Type)
	g.puts(")(")
	g.Exp(c.Exp)
	g.puts(")")
}

func (g *Gen) List(l ast.List) {
	if len(l.Exps) == 0 {
		g.puts("{}")
		return
	}

	g.puts("{")
	g.Exp(l.Exps[0])
	for _, e := range l.Exps[1:] {
		g.puts(", ")
		g.Exp(e)
	}
	g.puts("}")
}

func (g *Gen) Object(o ast.Object) {
	if len(o.Fields) == 0 {
		g.puts("{}")
		return
	}

	g.puts("{")
	g.nl()
	g.inc()

	for _, f := range o.Fields {
		g.indent()
		g.objField(f)
		g.puts(",")
		g.nl()
	}

	g.dec()
	g.indent()
	g.puts("}")
}

func (g *Gen) objField(f ast.ObjField) {
	g.puts(".")
	g.puts(f.Name.Str)
	g.puts(" = ")
	g.Exp(f.Exp)
}
