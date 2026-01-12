package stg

import (
	"fmt"

	"github.com/mebyus/ku/goku/compiler/enums/smk"
	"github.com/mebyus/ku/goku/compiler/enums/tpk"
)

type Map struct {
	mapkv

	// method list
	ml []*Symbol
}

type mapkv struct {
	Key   *Type
	Value *Type
}

// Explicit interface implementation check.
var _ TypeDef = &Map{}

func (*Map) Kind() tpk.Kind {
	return tpk.Map
}

func (m *Map) getMethod(name string) *Symbol {
	for _, s := range m.ml {
		if s.Name == name {
			return s
		}
	}
	return nil
}

func (x *TypeIndex) newMapType(kv mapkv) *Type {
	typ := &Type{
		// TODO: calculate size
		Kind: tpk.Map,
	}
	r := x.getRef(typ)

	const n = 3 // total number of methods
	ml := make([]*Symbol, 0, n)

	ml = append(ml, &Symbol{
		Name: "get",
		Kind: smk.MapMethod,
		Def: &Fun{
			Signature: Signature{
				Receiver: r,
				Params:   []*Type{kv.Key},
				Result:   x.getTuple([]*Type{kv.Value, x.Known.Bool}),
			},
		},
	})

	ml = append(ml, &Symbol{
		Name: "set",
		Kind: smk.MapMethod,
		Def: &Fun{
			Signature: Signature{
				Receiver: r,
				Params:   []*Type{kv.Key, kv.Value},
			},
		},
	})

	ml = append(ml, &Symbol{
		Name: "has",
		Kind: smk.MapMethod,
		Def: &Fun{
			Signature: Signature{
				Receiver: r,
				Params:   []*Type{kv.Key},
				Result:   x.Known.Bool,
			},
		},
	})

	if len(ml) != n {
		panic(fmt.Sprintf("inconsistent number of methods %d and %d", len(ml), n))
	}

	typ.Def = &Map{
		mapkv: kv,
		ml:    ml,
	}

	return typ
}
