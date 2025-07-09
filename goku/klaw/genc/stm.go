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
	case ast.ForRange:
		g.ForRange(s)
	case ast.If:
		g.If(s)
	case ast.JumpNext:
		g.JumpNext(s)
	case ast.JumpOut:
		g.JumpOut(s)
	case ast.Match:
		g.Match(s)
	case ast.Stub:
		g.Stub(s)
	case ast.Never:
		g.Never(s)
	case ast.Debug:
		g.Debug(s)
	case ast.Panic:
		g.Panic(s)
	case ast.Must:
		g.Must(s)
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
	g.level += 1

	for _, n := range b.Nodes {
		g.indent()
		g.Statement(n)
		g.nl()
	}

	g.level -= 1
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

func (g *Gen) ForRange(r ast.ForRange) {
	g.puts("for (")
	g.TypeSpec(r.Type)
	g.space()
	g.puts(r.Name.Str)
	g.puts(" = 0; ")
	g.puts(r.Name.Str)
	g.puts(" < ")
	g.Exp(r.Exp)
	g.puts("; ")
	g.puts(r.Name.Str)
	g.puts(" += 1) ")
	g.Block(r.Body)
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
		g.puts(" = {};")
		return
	}
	_, ok := v.Exp.(ast.Dirty)
	if ok {
		g.semi()
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

func (g *Gen) Match(m ast.Match) {
	g.puts("switch (")
	g.Exp(m.Exp)
	g.puts(") {")
	g.nl()

	for _, c := range m.Cases {
		g.MatchCase(c)
		g.nl()
	}
	g.MatchElseCase(m.Else)

	g.nl()
	g.indent()
	g.puts("}")
	g.nl()
}

func (g *Gen) MatchCase(c ast.MatchCase) {
	g.indent()
	g.puts("case ")
	g.Exp(c.List[0])
	g.puts(":")
	for _, exp := range c.List[1:] {
		g.nl()
		g.indent()
		g.puts("case ")
		g.Exp(exp)
		g.puts(":")
	}

	g.space()
	g.Block(c.Body)
	g.nl()
	g.indent()
	g.puts("break;")
	g.nl()
}

func (g *Gen) MatchElseCase(c *ast.Block) {
	if c == nil {
		return
	}

	g.indent()
	g.puts("default: ")
	g.Block(*c)
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

func (g *Gen) Panic(p ast.Panic) {
	g.puts("panic_pos(")
	g.str(p.Msg)
	g.puts(", ")
	g.textPosArgs(p.Pin)
	g.puts(");")
}

func (g *Gen) Must(m ast.Must) {
	g.puts("must_pos(")
	g.Exp(m.Exp)
	g.puts(", ")
	g.textPosArgs(m.Span().Pin)
	g.puts(");")
}

func (g *Gen) StaticMust(m ast.StaticMust) {
	g.puts("static_assert(")
	g.Exp(m.Exp)
	g.puts(");")
}

func (g *Gen) JumpNext(j ast.JumpNext) {
	g.puts("continue;")
}

func (g *Gen) JumpOut(j ast.JumpOut) {
	g.puts("break;")
}
