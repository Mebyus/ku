package parser

import (
	"github.com/mebyus/ku/internal/ku/ast"
	"github.com/mebyus/ku/internal/ku/token"
)

func (p *Parser) TypeSpec() (ast.TypeSpec, ss) {
	switch p.peek.Kind {
	case token.Word:
		pin := p.peek.Pin
		name := p.peek.Data
		p.advance() // skip word
		return &ast.TypeName{
			Name: name,
			Pin:  pin,
		}, 0
	default:
		pin := p.peek.Pin
		er := ast.Error{
			Pin: pin,
		}
		// error + sync
		return &ast.InvType{Error: er}, 0
	}
}
