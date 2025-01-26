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
