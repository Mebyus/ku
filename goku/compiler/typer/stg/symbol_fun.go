package stg

// Signature function or method signature.
type Signature struct {
	Params []*Type

	// Equals nil if function returns nothing or never returns.
	Result *Type

	// Always nil for regular functions.
	// Not nil only for methods.
	Receiver *Type

	// True for functions which never return.
	Never bool
}

// Fun symbol definition for functions and methods.
type Fun struct {
	nodeSymDef

	Body Block

	Signature
}
