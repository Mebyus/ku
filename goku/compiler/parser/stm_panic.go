package parser

import (
	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/token"
)

func (p *Parser) Panic() (ast.Panic, diag.Error) {
	pin := p.peek.Pin
	p.advance() // skip "panic"

	if p.peek.Kind != token.LeftParen {
		return ast.Panic{}, p.unexpected()
	}
	p.advance() // skip "("

	if p.peek.Kind != token.String {
		return ast.Panic{}, p.unexpected()
	}
	msg := p.peek.Data
	p.advance() // skip string

	if p.peek.Kind != token.RightParen {
		return ast.Panic{}, p.unexpected()
	}
	p.advance() // skip ")"

	if p.peek.Kind != token.Semicolon {
		return ast.Panic{}, p.unexpected()
	}
	p.advance() // skip ";"

	return ast.Panic{
		Pin: pin,
		Msg: msg,
	}, nil
}
