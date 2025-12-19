package parser

import "github.com/mebyus/ku/goku/compiler/diag"

func (p *Parser) unexpected() *diag.UnexpectedTokenError {
	return &diag.UnexpectedTokenError{Token: p.peek}
}
