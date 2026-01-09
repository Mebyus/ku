package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/stk"
	"github.com/mebyus/ku/goku/compiler/sm"
)

// DeferCall represents a statement with deferred function (or method) call.
type DeferCall struct {
	Call Call
}

var _ Statement = DeferCall{}

func (DeferCall) Kind() stk.Kind {
	return stk.DeferCall
}

func (c DeferCall) Span() sm.Span {
	return sm.Span{Pin: c.Call.Span().Pin}
}

func (c DeferCall) String() string {
	panic("not implemented")
	// var g Printer
	// g.Block(c)
	// return g.Output()
}
