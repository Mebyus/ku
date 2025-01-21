package parser

import (
	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/token"
)

func (p *Parser) Ret() (ast.Ret, diag.Error) {
	pin := p.c.Pin

	p.advance() // skip "ret"

	if p.c.Kind == token.Semicolon {
		p.advance() // skip ";"
		return ast.Ret{Pin: pin}, nil
	}

	exp, err := p.Pack()
	if err != nil {
		return ast.Ret{}, err
	}
	if p.c.Kind != token.Semicolon {
		return ast.Ret{}, p.unexpected()
	}
	p.advance() // skip ";"

	return ast.Ret{
		Pin: pin,
		Exp: exp,
	}, nil
}
