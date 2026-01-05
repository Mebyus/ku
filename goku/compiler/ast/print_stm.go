package ast

import (
	"fmt"
)

func (g *Printer) Statement(s Statement) {
	switch s := s.(type) {
	case nil:
		panic("nil statement")
	case Block:
		g.Block(s)
	case Ret:
		g.Ret(s)
	case Var:
		g.Var(s)
	case Const:
		g.Const(s)
	case Alias:
		g.Alias(s)
	case Assign:
		g.Assign(s)
	case Invoke:
		g.Invoke(s)
	case Loop:
		g.Loop(s)
	case While:
		g.While(s)
	case If:
		g.If(s)
	case JumpNext:
		g.JumpNext(s)
	case JumpOut:
		g.JumpOut(s)
	case Stub:
		g.Stub(s)
	case Never:
		g.Never(s)
	case Must:
		g.Must(s)
	case StaticMust:
		g.StaticMust(s)
	case Debug:
		g.Debug(s)
	default:
		panic(fmt.Sprintf("unexpected \"%s\" (=%d) statement (%T)", s.Kind(), s.Kind(), s))
	}
}

func (g *Printer) If(i If) {
	g.ifClause(i.If)
	for _, c := range i.ElseIfs {
		g.puts(" else ")
		g.ifClause(c)
	}
	if i.Else != nil {
		g.puts(" else ")
		g.Block(*i.Else)
	}
}

func (g *Printer) ifClause(c IfClause) {
	g.puts("if ")
	g.Exp(c.Exp)
	g.space()
	g.Block(c.Body)
}

func (g *Printer) Loop(l Loop) {
	g.puts("for ")
	g.Block(l.Body)
}

func (g *Printer) While(w While) {
	g.puts("for ")
	g.Exp(w.Exp)
	g.space()
	g.Block(w.Body)
}

func (g *Printer) ForRange(r ForRange) {
	g.puts("for ")
	g.puts(r.Name.Str)
	g.puts(": ")
	g.TypeSpec(r.Type)
	g.puts(" = [")
	if r.Start != nil {
		g.Exp(r.Start)
	}
	g.puts(":")
	g.Exp(r.End)
	g.puts("] ")
	g.Block(r.Body)
}

func (g *Printer) Stub(s Stub) {
	g.puts("stub;")
}

func (g *Printer) Never(n Never) {
	g.puts("never;")
}

func (g *Printer) Lookup(l Lookup) {
	g.puts("#lookup ")
	g.Exp(l.Exp)
	g.puts(";")
}

func (g *Printer) Static(s Static) {
	if len(s.Nodes) == 0 {
		g.puts("#{}")
		return
	}

	g.puts("#{")
	g.nl()
	g.inc()

	for _, n := range s.Nodes {
		g.indent()
		g.Statement(n)
		g.nl()
	}

	g.dec()
	g.indent()
	g.puts("}")
}

func (g *Printer) Block(b Block) {
	if len(b.Nodes) == 0 {
		g.puts("{}")
		return
	}

	g.puts("{")
	g.nl()
	g.inc()

	for _, n := range b.Nodes {
		g.indent()
		g.Statement(n)
		g.nl()
	}

	g.dec()
	g.indent()
	g.puts("}")
}

func (g *Printer) Debug(d Debug) {
	g.puts("#debug ")
	g.Block(d.Block)
}

func (g *Printer) Invoke(i Invoke) {
	g.Call(i.Call)
	g.semi()
}

func (g *Printer) Ret(r Ret) {
	if r.Exp == nil {
		g.puts("ret;")
		return
	}

	g.puts("ret ")
	g.Exp(r.Exp)
	g.semi()
}

func (g *Printer) TopConst(l TopConst) {
	g.Const(l.Const)
}

func (g *Printer) Const(l Const) {
	g.puts("const ")
	g.puts(l.Name.Str)
	g.puts(": ")
	g.TypeSpec(l.Type)
	g.puts(" = ")
	g.Exp(l.Exp)
	g.semi()
}

func (g *Printer) TopAlias(a TopAlias) {
	g.Alias(a.Alias)
}

func (g *Printer) Alias(a Alias) {
	g.puts("let ")
	g.puts(a.Name.Str)
	g.puts(" => ")
	g.Exp(a.Exp)
	g.semi()
}

func (g *Printer) TopVar(v TopVar) {
	g.Var(v.Var)
}

func (g *Printer) Var(v Var) {
	g.puts("var ")
	g.puts(v.Name.Str)
	g.puts(": ")
	g.TypeSpec(v.Type)

	if v.Exp == nil {
		g.semi()
		return
	}

	g.puts(" = ")
	g.Exp(v.Exp)
	g.semi()
}

func (g *Printer) Assign(a Assign) {
	g.Exp(a.Target)
	g.space()
	g.puts(a.Op.Kind.String())
	g.space()
	g.Exp(a.Value)
	g.semi()
}

func (g *Printer) Panic(p Panic) {
	g.puts("panic(\"")
	g.Exp(p.Exp)
	g.puts("\");")
}

func (g *Printer) Must(m Must) {
	g.puts("must(")
	g.Exp(m.Exp)
	g.puts(");")
}

func (g *Printer) Test(t Test) {
	g.puts("test(")
	g.Exp(t.Exp)
	g.puts(");")
}

func (g *Printer) StaticMust(m StaticMust) {
	g.puts("#must(")
	g.Exp(m.Exp)
	g.puts(");")
}

func (g *Printer) JumpNext(j JumpNext) {
	g.puts("jump @.next;")
}

func (g *Printer) JumpOut(j JumpOut) {
	g.puts("jump @.out;")
}
