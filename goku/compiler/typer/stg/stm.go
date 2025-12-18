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

type If struct {
	If     Branch
	IfElse []Branch
	Else   *Block
}

type Branch struct {
	Block Block

	// Branch condition. Always not nil.
	Exp Exp
}
