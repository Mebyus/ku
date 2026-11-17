package genc

import (
	"fmt"

	"github.com/mebyus/ku/internal/ku/enums/typk"
	"github.com/mebyus/ku/internal/ku/stg"
)

func (g *genbuf) typ(t *stg.Type) {
	g.puts(g.getTypeName(t))
}

func (g *genbuf) getTypeName(t *stg.Type) string {
	s, ok := g.types[t]
	if ok {
		return s
	}

	s = g.newTypeName(t)
	g.types[t] = s
	return s
}

func (g *genbuf) newTypeName(t *stg.Type) string {
	if t.Name != "" {
		if t.Name == "int" {
			return "ku_int"
		}
		if t.Unit == nil {
			return t.Name
		}
		return "t_" + t.Unit.LinkName + "_" + t.Name
	}

	switch t.Kind {
	case typk.Ref, typk.ArrayPointer:
		return g.getTypeName(t.Def.(*stg.Ref).Type) + "*"
	case typk.Span:
		return "s_" + g.getTypeName(t.Def.(*stg.Span).Type)
	default:
		panic(fmt.Sprintf("unexpected %s type", t.Kind))
	}
}
