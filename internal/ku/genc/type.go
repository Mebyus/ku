package genc

import (
	"github.com/mebyus/ku/internal/ku/enums/typk"
	"github.com/mebyus/ku/internal/ku/stg"
)

func (g *Buffer) typ(t *stg.Type) {
	g.puts(g.getTypeName(t))
}

func (g *Buffer) getTypeName(t *stg.Type) string {
	switch t.Kind {
	case typk.Integer:
		return "u32"
	default:
		return "invalid"
	}
}
