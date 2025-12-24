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
	// Always has at least 1 element.
	// First element is the first if branch.
	Branches []*Branch
	Else     *Block
}

type Branch struct {
	Block Block

	// Branch condition. Always not nil.
	Exp Exp

	Flags BranchFlag
}

type BranchFlag uint16

const (
	// Branch condition can be evaluated at compile-time.
	BranchStatic BranchFlag = 1 << iota

	// Only applicable if static flag is set.
	//
	// Branch condition is always true if flag is set.
	// Branch condition is always false if flag is not set.
	BranchTrue
)

func (b *Branch) SetFlags() {
	if b.Exp.Type().IsStatic() {
		b.Flags |= BranchStatic

		c := b.Exp.(*Boolean).Val
		if c {
			b.Flags |= BranchTrue
		}
	}
}

func (b *Branch) IsStatic() bool {
	return b.Flags&BranchStatic != 0
}

// IsTrue use only if branch is static.
func (b *Branch) IsTrue() bool {
	return b.Flags&BranchTrue != 0
}

// IsEmpty returns true if branch has no statements inside.
func (b *Branch) IsEmpty() bool {
	return len(b.Block.Nodes) == 0
}

// InvokeSymbol statement which directly (not via function pointer) calls
// a specific symbol (function or method).
type InvokeSymbol struct {
	Args []Exp

	Symbol *Symbol
}
