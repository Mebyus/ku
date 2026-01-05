package stg

// BuiltinGens contains builtin generic functions.
type BuiltinGens struct {
	Min *Symbol
	Max *Symbol
}

type GenIndex struct {
	Builtin BuiltinGens
}

func (g *GenIndex) Init() {}

