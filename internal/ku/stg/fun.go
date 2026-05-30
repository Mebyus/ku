package stg

// Signature represents function or method signature.
type Signature struct {
	Inputs []*Type

	// Equals nil if function returns nothing or never returns.
	Result *Type

	// Always nil for regular functions.
	// Not nil only for methods.
	Receiver *Type

	// True for functions which never return.
	Never bool
}

// FunDef symbol definition for functions and methods.
type FunDef struct {
	symdef

	Body Block

	Signature

	Params []*Symbol
}
