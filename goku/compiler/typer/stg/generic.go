package stg

type Generic struct {
	Scope Scope
}

var _ SymDef = &Generic{}

func (g *Generic) _symdef() {}
