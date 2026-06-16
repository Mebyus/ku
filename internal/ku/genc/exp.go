package genc

import (
	"fmt"

	"github.com/mebyus/ku/internal/ku/enums/bop"
	"github.com/mebyus/ku/internal/ku/stg"
)

func (g *Buffer) exp(exp stg.Exp) {
	switch e := exp.(type) {
	case *stg.Integer:
		if e.Neg {
			g.putb('-')
		}
		g.putn(e.Val)
	case *stg.Boolean:
		if e.Val {
			g.puts("true")
		} else {
			g.puts("false")
		}
	case *stg.SymExp:
		g.puts(g.getName(e.Symbol))
	case *stg.BinExp:
		g.binExp(e)
	case *stg.SpanNum:
		g.exp(e.Exp)
		g.puts(".num")
	default:
		panic(fmt.Sprintf("unexpected %T expression", e))
	}
}

func (g *Buffer) binExp(exp *stg.BinExp) {
	a, ok := exp.A.(*stg.BinExp)
	if ok {
		if shouldParen(a.Op.Kind, exp.Op.Kind) {
			g.putb('(')
			g.binExp(a)
			g.putb(')')
		} else {
			g.binExp(a)
		}
	} else {
		g.exp(exp.A)
	}

	// TODO: use C operator representation where it is required
	g.space()
	g.puts(exp.Op.Kind.String())
	g.space()

	b, ok := exp.B.(*stg.BinExp)
	if ok {
		if shouldParen(b.Op.Kind, exp.Op.Kind) {
			g.putb('(')
			g.binExp(b)
			g.putb(')')
		} else {
			g.binExp(b)
		}
	} else {
		g.exp(exp.B)
	}
}

// returns true if binary operator b would take precedence over a and
// that could lead to wrong order of operations inside binary expression
// if 3 or more operands
func shouldParen(a, b bop.Kind) bool {
	// TODO: use C precedence for operators
	ap := a.Precedence()
	bp := b.Precedence()

	if ap > bp {
		return true
	}

	// TODO: should check for operator associativity here instead of this hack-check
	return ap == bp && a != b
}
