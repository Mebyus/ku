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
	default:
		panic(fmt.Sprintf("unexpected %T expression", e))
	}
}
