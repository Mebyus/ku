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

	exp, err := p.Exp()
	if err != nil {
		return ast.Panic{}, err
	}

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
		Exp: exp,
	}, nil
}
