package ir

// Atom represents code atom inside function body.
// There two types of atoms:
//   - instruction
//   - label placement
type Atom interface {
	_atom()
}

type nodeAtom struct{}

func (nodeAtom) _atom() {}
