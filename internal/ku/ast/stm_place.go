package ast

import (
	"github.com/mebyus/ku/internal/fbuf"
	"github.com/mebyus/ku/internal/ku/sx"
)

// Place represents a pseudo-statement that places a label.
type Place struct {
	stm

	// label name
	Name string

	// label pin
	Pin sx.Pin
}

var _ Statement = &Place{}

func (p *Place) Format(buf *fbuf.Buffer) {
	buf.Puts("zone ")
	buf.Puts(p.Name)
	buf.Semi()
}
