package genc

import (
	"fmt"

	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/kub/eval"
)

func (g *Gen) evalBool(exp ast.Exp) bool {
	value, err := eval.EvalExp(g.State.Env, exp)
	if err != nil {
		panic(diag.Format(g.State.Map, err))
	}
	v, ok := value.(eval.Boolean)
	if !ok {
		panic(fmt.Sprintf("branch condition evaluates to (%T)", value))
	}

	return v.Val
}

func (g *Gen) StaticIf(s ast.StaticIf) {
	if g.evalBool(s.If.Exp) {
		g.Statements(s.If.Body.Nodes)
		return
	}

	for _, c := range s.ElseIfs {
		if g.evalBool(c.Exp) {
			g.Statements(c.Body.Nodes)
			return
		}
	}

	if s.Else != nil {
		g.Statements(s.Else.Nodes)
	}
}
