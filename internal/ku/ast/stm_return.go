package ast

import (
	"github.com/mebyus/ku/internal/fbuf"
	"github.com/mebyus/ku/internal/ku/sx"
)

// Return represents return statement.
//
// Formal definition:
//
//	Return -> "return" [ Exp ] ";"
type Return struct {
	stm

	// Equals nil if return does not have expression.
	Exp Exp

	// Return keyword pin.
	Pin sx.Pin
}

var _ Statement = &Return{}

func (r *Return) Format(buf *fbuf.Buffer) {
	if r.Exp == nil {
		buf.Puts("return;")
		return
	}

	buf.Puts("return ")
	r.Exp.Format(buf)
	buf.Semi()
}
