package ast

import (
	"io"
)

func Render(w io.Writer, text *Text) error {
	var p Printer
	p.Text(text)
	_, err := p.WriteTo(w)
	return err
}

func (g *Printer) Text(text *Text) {

}
