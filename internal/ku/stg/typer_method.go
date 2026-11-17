package stg

import (
	"fmt"

	"github.com/mebyus/ku/internal/ku/ast"
	"github.com/mebyus/ku/internal/ku/enums/symk"
)

// bind method receivers to their custom types
func (t *Typer) bind() {
	for i, s := range t.methods {
		m := t.box.methods[i]
		name := m.Recv.Name
		pin := m.Recv.Pin
		rtyp := m.Recv.Type

		r := t.unit.Scope.Get(name)
		if r == nil {
			t.report(pin, 0, fmt.Sprintf("method receiver \"%s\" refers to undefined symbol", name))
			continue
		}
		if r.Kind != symk.Type {
			t.report(pin, 0, fmt.Sprintf("method receiver \"%s\" refers to %s symbol (instead of custom type)", name, r.Kind))
			continue
		}

		t.bindMethod(name, rtyp, r, s)
	}
}

func (t *Typer) bindMethod(name string, rtyp uint32, r, s *Symbol) {
	typ := r.Def.(*Type)
	c := typ.Def.(*Custom)
	c.bind(name, s)

	var rt *Type
	switch rtyp {
	case 0:
		panic("not a method")
	case ast.RecvVal:
		rt = typ
	case ast.RecvRef:
		rt = t.common.Types.getRef(typ)
	case ast.RecvPtr:
		panic("TODO")
	default:
		panic(fmt.Sprintf("invalid receiver type (=%d)", rtyp))
	}

	s.Def.(*FunDef).Receiver = rt
}
