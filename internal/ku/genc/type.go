package genc

import (
	"fmt"

	"github.com/mebyus/ku/internal/ku/enums/typk"
	"github.com/mebyus/ku/internal/ku/stg"
)

func (g *Buffer) typ(t *stg.Type) {
	g.puts(g.getTypeName(t))
}

func (g *Buffer) getTypeName(t *stg.Type) string {
	s, ok := g.types[t]
	if ok {
		return s
	}

	s = g.newTypeName(t)
	g.types[t] = s
	return s
}

func (g *Buffer) newTypeName(t *stg.Type) string {
	switch t.Kind {
	case typk.Integer:
		if t.IsSigned() {
			switch t.Size {
			case 1:
				return "s8"
			case 2:
				return "s16"
			case 4:
				return "s32"
			case 8:
				return "s64"
			default:
				panic(fmt.Sprintf("unexpected integer size %d", t.Size))
			}
		}
		switch t.Size {
		case 1:
			return "u8"
		case 2:
			return "u16"
		case 4:
			return "u32"
		case 8:
			return "u64"
		default:
			panic(fmt.Sprintf("unexpected integer size %d", t.Size))
		}
	case typk.Boolean:
		return "bool"
	default:
		panic(fmt.Sprintf("unexpected %s type", t.Kind))
	}
}
