package genc

import (
	"fmt"

	"github.com/mebyus/ku/internal/ku/enums/typk"
	"github.com/mebyus/ku/internal/ku/stg"
)

// produce a type specifier suitable for usage in places like:
//   - variable type declaration
//   - function return type
//   - type of function parameter
func (g *genbuf) typespec(t *stg.Type) {
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
		return g.mangleTypeName(t)
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

func (g *genbuf) typedef(t *stg.Symbol) {
	g.puts("typedef ")
	g.typ(t.Def.(*stg.Type).Def.(*stg.Custom).Type)
	g.space()
	g.puts(g.getName(t))
	g.semi()
	g.line()
}

// produce type definition suitable for C typedef
func (g *genbuf) typ(t *stg.Type) {
	if t.Name != "" {
		g.puts(g.getTypeName(t))
		return
	}

	panic("TODO")
}

// get C typename for a type that has non-empty Type.Name
func (g *genbuf) mangleTypeName(t *stg.Type) string {
	if t.Name == "int" {
		return "ku_int"
	}
	if t.Unit == nil {
		return t.Name
	}
	return t.Unit.LinkName + "_t_" + t.Name
}
