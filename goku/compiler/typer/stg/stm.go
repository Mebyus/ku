package stg

import "github.com/mebyus/ku/goku/compiler/sm"

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

	// Assign target. Always not nil.
	Target Exp
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

// Invoke statement which call something without assigning the result.
type Invoke struct {
	// Must be one of call expressions.
	Call Exp
}

// While represents conditional loop statement.
type While struct {
	Body Block

	// Loop condition. Always not nil.
	Exp Exp
}

// Loop represents unconditional loop statement.
type Loop struct {
	Body Block
}

// Must represents must statement.
type Must struct {
	// Always not nil.
	Exp Exp
}

// Stub represents stub statement.
type Stub struct {
	Pin sm.Pin
}

// Never represents never statement.
type Never struct {
	Pin sm.Pin
}

// ForRange represents for range statement.
type ForRange struct {
	Body Block

	// Equals nil if omitted. Zero value is used.
	Start Exp

	// Always not nil. Always has integer type.
	End Exp

	// Loop variable.
	Var *Symbol
}

type MatchInteger struct {
	Cases []*MatchCase

	// Expression being matched. Always has integer or enum type.
	Exp Exp

	Else *Block
}

type MatchCase struct {
	Body Block

	// Always has at least one element.
	List []Exp
}
