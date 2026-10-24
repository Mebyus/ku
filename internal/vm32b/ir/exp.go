package ir

import "github.com/mebyus/ku/internal/ku/enums/bop"

type Exp interface {
	_exp()
}

type exp struct{}

func (exp) _exp() {}

// Interger constant used as expression.
type Integer struct {
	exp

	Val uint32
}

var _ Exp = &Integer{}

// local variable used as expression
type Local struct {
	exp

	Name uint32
}

var _ Exp = &Local{}

// take address of local variable
type Addr struct {
	exp

	Name uint32
}

var _ Exp = &Addr{}

// binary expression
type BinExp struct {
	exp

	A, B Exp
	Op   bop.Kind
}

var _ Exp = &BinExp{}

type Call struct {
	exp

	Args []Exp

	// function name (=index)
	Name uint32
}

var _ Exp = &Call{}

// slice local array
type SliceLocal struct {
	exp

	// start index
	A Exp

	// end index
	B Exp

	// local slot name, where to slice
	Name uint32

	// number of slots occupied by one element
	// cannot be 0
	Stride uint32
}

var _ Exp = &SliceLocal{}

// take element by index starting from address
type Index struct {
	exp

	Addr Exp

	Index Exp

	// number of slots occupied by one element
	// cannot be 0
	Stride uint32
}

var _ Exp = &Index{}
