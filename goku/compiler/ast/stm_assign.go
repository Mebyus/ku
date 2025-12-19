package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/aok"
	"github.com/mebyus/ku/goku/compiler/enums/stk"
	"github.com/mebyus/ku/goku/compiler/sm"
)

// Assign represents an assign statement.
//
// Formal definition:
//
//	Assign => Target AssignOp Value ";"
//	Target => Exp
//	Value  => Exp
type Assign struct {
	// Target of assignment operation. The value referred to by this expression
	// is altered by the statement.
	Target Exp

	// Value which is used to perform the alteration of target.
	Value Exp

	// Specifies which kind of assignment this operation is.
	Op AssignOp
}

// Explicit interface implementation check.
var _ Statement = Assign{}

func (Assign) Kind() stk.Kind {
	return stk.Assign
}

func (a Assign) Span() sm.Span {
	return sm.Span{Pin: a.Op.Pin}
}

func (a Assign) String() string {
	var g Printer
	g.Assign(a)
	return g.Output()
}

// AssignOp represents assign operator inside expression.
type AssignOp struct {
	Pin  sm.Pin
	Kind aok.Kind
}
