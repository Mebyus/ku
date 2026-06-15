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
	case token.LeftSquare:
		if p.next.Kind == token.RightSquare {
			return p.typeSpan()
		}
	}

	pin := p.peek.Pin
	er := ast.Error{
		Pin: pin,
	}
	// error + sync
	return &ast.InvType{Error: er}, 0
}

func (p *Parser) typeSpan() (ast.TypeSpec, ss) {
	pin := p.peek.Pin

	p.advance() // skip "["
	p.advance() // skip "]"

	typ, s := p.TypeSpec()
	if s != 0 {
		return &ast.InvType{Error: ast.Error{Pin: pin}}, s
	}

	return &ast.Span{
		Pin:  pin,
		Type: typ,
	}, 0
}
