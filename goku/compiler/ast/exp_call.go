package ast

import "github.com/mebyus/ku/goku/compiler/enums/exk"

// Call represents an expression of calling something.
//
//	Call => Chain "(" Args ")"
//	Args => { Exp "," } // trailing comma is optional
type Call struct {
	Chain Chain
	Args  []Exp
}

var _ Exp = Call{}

func (Call) Kind() exk.Kind {
	return exk.Call
}
