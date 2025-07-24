package parser

import (
	"fmt"

	"github.com/mebyus/ku/goku/compiler/diag"
)

func (p *Parser) unexpected() *diag.SimpleMessageError {
	return &diag.SimpleMessageError{
		Pin:  p.peek.Pin,
		Text: fmt.Sprintf("unexpected token %s", p.peek),
	}
}
