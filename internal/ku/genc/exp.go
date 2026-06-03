package genc

import (
	"fmt"

	"github.com/mebyus/ku/internal/ku/stg"
)

func (g *Buffer) exp(exp stg.Exp) {
	switch e := exp.(type) {
	case *stg.Integer:
		if e.Neg {
			g.putb('-')
		}
		g.putn(e.Val)
	case *stg.BinExp:
		g.exp(e.A)
		g.space()
		g.puts(e.Op.Kind.String())
		g.space()
		g.exp(e.B)
	default:
		panic(fmt.Sprintf("unexpected %T expression", e))
	}
}
