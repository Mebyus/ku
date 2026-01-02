package stg

import (
	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/sm"
)

type UnaryOp = ast.UnaryOp

// Unary represents unary expression.
type Unary struct {
	Exp Exp
	Op  UnaryOp
}

func (u *Unary) Type() *Type {
	return u.Exp.Type()
}

func (u *Unary) Span() sm.Span {
	return sm.Span{Pin: u.Op.Pin}
}

func (u *Unary) String() string {
	panic("not implemented")
}
