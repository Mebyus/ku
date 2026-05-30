package parser

import (
	"github.com/mebyus/ku/internal/ku/ast"
	"github.com/mebyus/ku/internal/ku/token"
)

func (p *Parser) TypeSpec() (ast.TypeSpec, bool) {
	switch p.peek.Kind {
	case token.Word:
		pin := p.peek.Pin
		name := p.peek.Data
		p.advance() // skip word
		return &ast.TypeName{
			Name: name,
			Pin:  pin,
		}, true
	default:
		panic("stub")
	}
}
