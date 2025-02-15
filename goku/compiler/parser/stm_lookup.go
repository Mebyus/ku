package parser

import (
	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/token"
)

func (p *Parser) Lookup() (ast.Lookup, diag.Error) {
	pin := p.c.Pin
	p.advance() // skip "#lookup"

	if p.c.Kind != token.LeftParen {
		return ast.Lookup{}, p.unexpected()
	}
	p.advance() // skip "("

	exp, err := p.Exp()
	if err != nil {
		return ast.Lookup{}, err
	}

	if p.c.Kind != token.RightParen {
		return ast.Lookup{}, p.unexpected()
	}
	p.advance() // skip ")"

	if p.c.Kind != token.Semicolon {
		return ast.Lookup{}, p.unexpected()
	}
	p.advance() // skip ";"

	return ast.Lookup{
		Pin: pin,
		Exp: exp,
	}, nil
}
