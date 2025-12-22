package stg

type Statement interface{}

// Block represents block statement or function body.
type Block struct {
	Scope Scope

	Nodes []Statement
}

// Ret represents return statement.
type Ret struct {
	// Can be nil, if return does not have expression.
	Exp Exp
}

// Var represents var declaration statement.
type Var struct {
	Symbol *Symbol

	// Initial expression. Always not nil.
	Exp Exp
}

// Assign represents value assignment to a variable.
type Assign struct {
	// Always not nil.
	Exp Exp

	Symbol *Symbol
}

type If struct {
	If      Branch
	ElseIfs []Branch
	Else    *Block
}

type Branch struct {
	Block Block

	// Branch condition. Always not nil.
	Exp Exp
}

// InvokeSymbol statement which directly (not via function pointer) calls
// a specific symbol (function or method).
type InvokeSymbol struct {
	Args []Exp

	Symbol *Symbol
}
