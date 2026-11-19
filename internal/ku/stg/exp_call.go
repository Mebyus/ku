package stg

import "github.com/mebyus/ku/internal/ku/sx"

// FunCall represents expression which calls a specific function.
type FunCall struct {
	operand

	// can be empty, if call has no arguments
	Args []Exp

	// pin of call site
	pin sx.Pin

	// What function is being called.
	Symbol *Symbol

	// Call result type. Always not nil.
	typ *Type
}

// Explicit interface implementation check.
var _ Operand = &FunCall{}

func (c *FunCall) Type() *Type {
	return c.typ
}

func (c *FunCall) Pin() sx.Pin {
	return c.pin
}

func (c *FunCall) Target() bool {
	return false
}

// MethodCall represents expression which calls a specific method.
type MethodCall struct {
	operand

	// Call arguments, except method receiver.
	// Can be empty if call does not have arguments.
	Args []Exp

	// Receiver value for this call.
	// Always not nil.
	Recv Exp

	// where call is made
	pin sx.Pin

	// What method is being called.
	Symbol *Symbol

	// Call result type. Always not nil.
	typ *Type
}

var _ Operand = &MethodCall{}

func (c *MethodCall) Type() *Type {
	return c.typ
}

func (c *MethodCall) Pin() sx.Pin {
	return c.pin
}

func (c *MethodCall) Target() bool {
	return false
}
