package ast

import (
	"github.com/mebyus/ku/internal/fbuf"
	"github.com/mebyus/ku/internal/ku/sx"
)

// InvokeSymbol represents statement formed by function-like
// call to a symbol.
type InvokeSymbol struct {
	stm

	// can be nil if call does not have arguments
	Args []Exp

	// name of the called symbol
	Name string

	// name pin
	Pin sx.Pin
}

var _ Statement = &InvokeSymbol{}

func (s *InvokeSymbol) Format(buf *fbuf.Buffer) {
	buf.Puts(s.Name)
	if len(s.Args) == 0 {
		buf.Puts("();")
		return
	}

	buf.Putb('(')
	s.Args[0].Format(buf)
	for _, arg := range s.Args[1:] {
		buf.Puts(", ")
		arg.Format(buf)
	}
	buf.Puts(");")
}
