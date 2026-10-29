package stg

import "github.com/mebyus/ku/internal/ku/sx"

// IfRet represents short if-return statement.
type IfRet struct {
	stm

	// condition, must have boolean type
	//
	// always not nil
	Exp Exp

	// true branch
	//
	// may be nil if return does not have expression
	Ret Exp

	// "if" pin
	pin sx.Pin
}

var _ Statement = &IfRet{}

func (f *IfRet) Pin() sx.Pin {
	return f.pin
}
