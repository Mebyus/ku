package ast

import (
	"github.com/mebyus/ku/internal/fbuf"
	"github.com/mebyus/ku/internal/ku/sx"
)

// Const represents constant definition statement.
//
// Formal definition:
//
//	Const -> "const" Name ":" [ TypeSpec ] "=" Exp ";"
type Const struct {
	stm

	// Can be nil if constant type is not specified.
	Type TypeSpec

	// Specifies constant init value expression.
	//
	// Always not nil.
	Exp Exp

	Name string

	// Constant name pin.
	Pin sx.Pin
}

var _ Statement = &Const{}

func (c *Const) Format(buf *fbuf.Buffer) {
	buf.Puts("const ")
	buf.Puts(c.Name)

	if c.Type != nil {
		buf.Puts(": ")
		c.Type.Format(buf)
		buf.Puts(" = ")
	} else {
		buf.Puts(" := ")
	}

	c.Exp.Format(buf)
	buf.Semi()
}
