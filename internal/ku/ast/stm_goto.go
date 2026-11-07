package ast

import (
	"github.com/mebyus/ku/internal/fbuf"
	"github.com/mebyus/ku/internal/ku/sx"
)

// Goto represents goto label statement.
type Goto struct {
	stm

	// label name
	Name string

	// lavel pin
	Pin sx.Pin
}

var _ Statement = &Goto{}

func (g *Goto) Format(buf *fbuf.Buffer) {
	buf.Puts("goto ")
	buf.Puts(g.Name)
	buf.Semi()
}
