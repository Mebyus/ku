package ast

import (
	"fmt"

	"github.com/mebyus/ku/goku/compiler/char"
)

func (g *Printer) Exp(exp Exp) {
	switch e := exp.(type) {
	case nil:
		panic("nil exp")
	case Symbol:
		g.Symbol(e)
	case DotName:
		g.DotName(e)
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
	case GetRef:
		g.GetRef(e)
	case Slice:
		g.Slice(e)
	case Tweak:
		g.Tweak(e)
	case Object:
		g.Object(e)
	case List:
		g.List(e)
	case TypeId:
		g.TypeId(e)
	case ErrorId:
		g.ErrorId(e)
	default:
		panic(fmt.Sprintf("unexpected \"%s\" (=%d) expression (%T)", e.Kind(), e.Kind(), e))
	}
}

func (g *Printer) Tweak(t Tweak) {
	g.Chain(t.Chain)
	g.tweakFields(t.Fields)
}

func (g *Printer) tweakFields(fields []ObjField) {
	if len(fields) == 0 {
		g.puts(".{}")
		return
	}

	g.puts(".{")
	g.inc()
	for _, f := range fields {
		g.nl()
		g.indent()
		g.ObjField(f)
	}
	g.dec()
	g.nl()
	g.indent()
	g.puts("}")
}

func (g *Printer) Object(o Object) {
	if len(o.Fields) == 0 {
		g.puts("{}")
		return
	}

	g.puts("{")
	g.inc()
	for _, f := range o.Fields {
		g.nl()
		g.indent()
		g.ObjField(f)
	}
	g.dec()
	g.nl()
	g.indent()
	g.puts("}")
}

func (g *Printer) ObjField(f ObjField) {
	g.puts(f.Name.Str)
	g.puts(": ")
	g.Exp(f.Exp)
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

func (g *Printer) GetRef(r GetRef) {
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
	if len(c.Parts) == 0 {
		return
	}
	if c.Start.Str == "" {
		p := c.Parts[0]
		u, ok := p.(Unsafe)
		if ok {
			g.puts("unsafe.")
			g.puts(u.Name)
		} else {
			g.Part(p)
		}
	} else {
		g.Part(c.Parts[0])
	}
	for _, p := range c.Parts[1:] {
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
	case Unsafe:
		g.Unsafe(p)
	default:
		panic(fmt.Sprintf("unexpected \"%s\" (=%d) chain part (%T)", p.Kind(), p.Kind(), p))
	}
}

func (g *Printer) SelectTest(s SelectTest) {
	g.puts(".test.")
	g.puts(s.Name.Str)
}

func (g *Printer) Deref(d Deref) {
	g.puts(".*")
}

func (g *Printer) DerefSelect(d DerefSelect) {
	g.puts(".*.")
	g.puts(d.Name.Str)
}

func (g *Printer) Select(s Select) {
	g.puts(".")
	g.puts(s.Name.Str)
}

func (g *Printer) Unsafe(u Unsafe) {
	g.puts(".unsafe.")
	g.puts(u.Name)
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

func (g *Printer) List(l List) {
	if len(l.Exps) == 0 {
		g.puts("[]")
		return
	}

	g.puts("[")
	g.Exp(l.Exps[0])
	for _, e := range l.Exps[1:] {
		g.puts(", ")
		g.Exp(e)
	}
	g.puts("]")
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

func (g *Printer) Size(s Size) {
	g.puts("#size(")
	g.TypeSpec(s.Exp)
	g.puts(")")
}

func (g *Printer) Cast(c Cast) {
	g.puts("#cast(")
	g.TypeSpec(c.Type)
	g.puts(", ")
	g.Exp(c.Exp)
	g.puts(")")
}

func (g *Printer) CheckFlag(c CheckFlag) {
	g.puts("#check(")
	g.Exp(c.Exp)
	g.puts(", ")
	g.Exp(c.Flag)
	g.puts(")")
}

func (g *Printer) ArrayLen(l ArrayLen) {
	g.puts("#len(")
	g.Exp(l.Exp)
	g.puts(")")
}

func (g *Printer) Tint(t Tint) {
	g.puts("tint(")
	g.TypeSpec(t.Type)
	g.puts(", ")
	g.Exp(t.Exp)
	g.puts(")")
}

func (g *Printer) Symbol(n Symbol) {
	g.puts(n.Name)
}

func (g *Printer) DotName(d DotName) {
	g.puts(".")
	g.puts(d.Name)
}

func (g *Printer) Integer(n Integer) {
	g.puts(n.String())
}

func (g *Printer) String(s String) {
	g.puts("\"")
	g.puts(char.Escape(s.Val))
	g.puts("\"")
}

func (g *Printer) Rune(s Rune) {
	g.puts("'")
	g.puts(char.EscapeRune(rune(s.Val)))
	g.puts("'")
}

func (g *Printer) Nil(n Nil) {
	g.puts("nil")
}

func (g *Printer) Dirty(d Dirty) {
	g.puts("?")
}

func (g *Printer) TypeId(t TypeId) {
	g.puts("#typeid(")
	g.puts(t.Name.Str)
	g.puts(")")
}

func (g *Printer) ErrorId(e ErrorId) {
	g.puts("#error(")
	g.puts(e.Name.Str)
	g.puts(")")
}

func (g *Printer) EnumMacro(e EnumMacro) {
	g.puts("#enum(")
	g.puts(e.Name.Str)
	g.puts(".")
	g.puts(e.Entry.Str)
	g.puts(")")
}

func (g *Printer) BuildQuery(q BuildQuery) {
	g.puts("#build.")
	g.puts(q.Name)
}
