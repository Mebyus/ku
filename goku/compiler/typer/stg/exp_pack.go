package stg

import (
	"strings"

	"github.com/mebyus/ku/goku/compiler/sm"
)

// Pack represents a pack expression.
type Pack struct {
	// Always contains at least 2 elements.
	List []Exp

	typ *Type
}

// Explicit interface implementation check.
var _ Exp = &Pack{}

func (p *Pack) Type() *Type {
	return p.typ
}

func (p *Pack) Span() sm.Span {
	return sm.Span{Pin: p.List[0].Span().Pin}
}

func (p *Pack) String() string {
	var g strings.Builder

	g.WriteString("(")
	g.WriteString(p.List[0].String())
	for _, e := range p.List[1:] {
		g.WriteString(", ")
		g.WriteString(e.String())
	}
	g.WriteString(")")

	return g.String()
}

func (x *TypeIndex) MakePack(list []Exp) *Pack {
	types := make([]*Type, 0, len(list))
	for _, e := range list {
		types = append(types, e.Type())
	}
	typ := x.getTuple(types)

	return &Pack{
		List: list,
		typ:  typ,
	}
}
