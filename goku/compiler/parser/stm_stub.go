package parser

import (
	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/token"
)

func (p *Parser) Stub() (ast.Stub, diag.Error) {
	pin := p.c.Pin

	p.advance() // skip "#stub"

	if p.c.Kind != token.Semicolon {
		return ast.Stub{}, p.unexpected()
	}
	p.advance() // skip ";"

	return ast.Stub{Pin: pin}, nil
}

func (p *Parser) Never() (ast.Never, diag.Error) {
	pin := p.c.Pin

	p.advance() // skip "#never"

	if p.c.Kind != token.Semicolon {
		return ast.Never{}, p.unexpected()
	}
	p.advance() // skip ";"

	return ast.Never{Pin: pin}, nil
}
