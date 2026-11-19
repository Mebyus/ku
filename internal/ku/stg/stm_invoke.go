package stg

import "github.com/mebyus/ku/internal/ku/sx"

// InvokeFun represents statement which calls a specific function.
type InvokeFun struct {
	stm

	// can be empty, if call has no arguments
	Args []Exp

	// pin of call site
	pin sx.Pin

	Symbol *Symbol
}

var _ Statement = &InvokeFun{}

func (v *InvokeFun) Pin() sx.Pin {
	return v.pin
}

// InvokeMethod represents statement which calls a specific method.
type InvokeMethod struct {
	stm

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
}

var _ Statement = &InvokeMethod{}

func (v *InvokeMethod) Pin() sx.Pin {
	return v.pin
}
