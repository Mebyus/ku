package parser

import (
	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/token"
)

func (p *Parser) Lookup() (ast.Lookup, diag.Error) {
	pin := p.peek.Pin
	p.advance() // skip "#lookup"

	exp, err := p.Exp()
	if err != nil {
		return ast.Lookup{}, err
	}

	if p.peek.Kind != token.Semicolon {
		return ast.Lookup{}, p.unexpected()
	}
	p.advance() // skip ";"

	return ast.Lookup{
		Pin: pin,
		Exp: exp,
	}, nil
}
