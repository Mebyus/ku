package genc

import (
	"fmt"

	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/enums/aok"
)

func (g *Gen) Statement(s ast.Statement) {
	switch s := s.(type) {
	case nil:
		panic("nil statement")
	case ast.Block:
		g.Block(s)
	case ast.Ret:
		g.Ret(s)
	case ast.Var:
		g.Var(s)
	case ast.Const:
		g.Const(s)
	case ast.Alias:
		panic("not supported")
	case ast.Assign:
		g.Assign(s)
	case ast.Invoke:
		g.Invoke(s)
	case ast.Loop:
		g.Loop(s)
	case ast.While:
		g.While(s)
	case ast.If:
		g.If(s)
	case ast.Stub:
		g.Stub(s)
	case ast.Never:
		g.Never(s)
	case ast.Debug:
		g.Debug(s)
	case ast.StaticMust:
		g.StaticMust(s)
	default:
		panic(fmt.Sprintf("unexpected \"%s\" (=%d) statement (%T)", s.Kind(), s.Kind(), s))
	}
}

func (g *Gen) Block(b ast.Block) {
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

func (g *Gen) Debug(d ast.Debug) {
	if !g.debug() {
		return
	}

	g.Block(d.Block)
}

func (g *Gen) If(i ast.If) {
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

func (g *Gen) ifClause(c ast.IfClause) {
	g.puts("if (")
	g.Exp(c.Exp)
	g.puts(") ")
	g.Block(c.Body)
}

func (g *Gen) Never(n ast.Never) {
	g.puts("panic_never(")
	g.textPosArgs(n.Pin)
	g.puts(");")
}

func (g *Gen) Loop(l ast.Loop) {
	g.puts("while (true) ")
	g.Block(l.Body)
}

func (g *Gen) While(w ast.While) {
	g.puts("while (")
	g.Exp(w.Exp)
	g.puts(") ")
	g.Block(w.Body)
}

func (g *Gen) TopVar(v ast.TopVar) {
	g.puts("static ")
	g.NameDef(v.Name.Str, v.Type)

	if v.Exp == nil {
		g.semi()
		return
	}

	g.puts(" = ")
	g.Exp(v.Exp)
	g.semi()
}

func (g *Gen) Var(v ast.Var) {
	g.NameDef(v.Name.Str, v.Type)

	if v.Exp == nil {
		g.puts("{};")
		return
	}

	g.puts(" = ")
	g.Exp(v.Exp)
	g.semi()
}

func (g *Gen) Const(c ast.Const) {
	g.puts("const ")
	g.NameDef(c.Name.Str, c.Type)
	g.puts(" = ")
	g.Exp(c.Exp)
	g.semi()
}

func (g *Gen) TopConst(l ast.TopConst) {
	g.puts("static ")
	g.Const(l.Const)
}

func (g *Gen) Assign(a ast.Assign) {
	if a.Op.Kind == aok.Simple {
		_, ok := a.Target.(ast.Blank)
		if ok {
			g.puts("(void)(")
			g.Exp(a.Value)
			g.puts(");")
			return
		}
	}

	g.Exp(a.Target)
	g.space()
	g.puts(a.Op.Kind.String())
	g.space()
	g.Exp(a.Value)
	g.semi()
}

func (g *Gen) Invoke(i ast.Invoke) {
	g.Call(i.Call)
	g.semi()
}

func (g *Gen) Ret(r ast.Ret) {
	if r.Exp == nil {
		g.puts("return;")
		return
	}

	g.puts("return ")
	g.Exp(r.Exp)
	g.semi()
}

func (g *Gen) Stub(s ast.Stub) {
	g.puts("panic_stub(")
	g.textPosArgs(s.Pin)
	g.puts(");")
}

func (g *Gen) StaticMust(m ast.StaticMust) {
	g.puts("static_assert(")
	g.Exp(m.Exp)
	g.puts(");")
}
