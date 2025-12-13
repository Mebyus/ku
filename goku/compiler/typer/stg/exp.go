package stg

// Exp node that represents an arbitrary expression.
type Exp interface {
	Type() *Type

	_symdef()
}
