package stg

import (
	"github.com/mebyus/ku/goku/compiler/enums/bgk"
	"github.com/mebyus/ku/goku/compiler/enums/smk"
)

type GenIndex struct {
	// Contains instantiated "min" functions for each specific configuration of
	// number of params + their type.
	//
	// Maps min function spec to corresponding symbol.
	Mins map[UniformParamsSpec]*Symbol
}

// UniformParamsSpec defines signature for builtin generic functions with uniform,
// but arbitrary number of params, such as "min" or "max".
type UniformParamsSpec struct {
	nodeSymDef

	// Param type.
	Type *Type

	// Number of params. Always not 0.
	Num uint
}

func (g *GenIndex) Init() {
	g.Mins = make(map[UniformParamsSpec]*Symbol)
}

func (g *GenIndex) getMinInstance(spec UniformParamsSpec) *Symbol {
	if spec.Num == 0 {
		panic("zero param number")
	}
	if spec.Type == nil {
		panic("param type not specified")
	}

	s := g.Mins[spec]
	if s != nil {
		return s
	}

	s = &Symbol{
		Name: "min",
		Kind: smk.BgenFunInst,
		Aux:  uint32(bgk.Min),
		Def:  spec,
	}

	g.Mins[spec] = s
	return s
}
