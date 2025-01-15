package ast

import "fmt"

func (g *Printer) Statement(s Statement) {
	switch s := s.(type) {
	case Block:
		g.Block(s)
	case Ret:
		g.Ret(s)
	default:
		panic(fmt.Sprintf("unexpected \"%s\" (=%d) statement (%T)", s.Kind(), s.Kind(), s))
	}
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

}

func (g *Printer) Let(l Let) {

}

func (g *Printer) TopVar(l TopVar) {

}

func (g *Printer) Var(l Var) {

}
