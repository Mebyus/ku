package ast

// Atom represents code atom inside function body.
// There two types of atoms:
//   - instruction
//   - label placement
type Atom interface {
	_atom()
}

type nodeAtom struct{}

func (nodeAtom) _atom() {}

// Operand represents instruction operand of any type.
type Operand interface {
	_operand()
}

type nodeOperand struct{}

func (nodeOperand) _operand() {}
