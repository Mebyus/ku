package stg

import (
	"github.com/mebyus/ku/internal/ku/ast"
	"github.com/mebyus/ku/internal/ku/enums/typk"
)

// evaluate and layout a custom type from its type specifier
func (t *Typer) evalTypeDef(s *Symbol, spec ast.TypeSpec) {
	ct := s.Def.(*Type)

	typ := t.LookupType(&t.unit.Scope, spec)
	if typ.Kind == typk.Custom {
		// we do not want to form deep chains of nesting
		// custom types within each other
		typ = typ.Def.(*Custom).Type
	}

	ct.Def.(*Custom).Type = typ
	ct.Size = typ.Size
	ct.Align = typ.Align
}
