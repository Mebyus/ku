package ast

import "fmt"

func (g *Printer) Statement(s Statement) {
	switch s := s.(type) {
	case Block:
		g.Block(s)
	case Ret:
		g.Ret(s)
	case Var:
		g.Var(s)
	case Let:
		g.Let(s)
	case Assign:
		g.Assign(s)
	case Invoke:
		g.Invoke(s)
	case Loop:
		g.Loop(s)
	default:
		panic(fmt.Sprintf("unexpected \"%s\" (=%d) statement (%T)", s.Kind(), s.Kind(), s))
	}
}

func (g *Printer) Loop(l Loop) {
	g.puts("for ")
	g.Block(l.Body)
}

func (g *Printer) Block(b Block) {
	if len(b.Nodes) == 0 {
		g.puts("{}")
		return
	}

	g.puts("{")
	g.nl()
	g.inc()

	for _, s := range b.Nodes {
		g.indent()
		g.Statement(s)
		g.nl()
	}

	g.dec()
	g.indent()
	g.puts("}")
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

func (g *Printer) TopLet(l TopLet) {
	g.Let(l.Let)
}

func (g *Printer) Let(l Let) {
	g.puts("let ")
	g.puts(l.Name.Str)
	g.puts(": ")
	g.TypeSpec(l.Type)
	g.puts(" = ")
	g.Exp(l.Exp)
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
