package parser

import (
	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/token"
)

func (p *Parser) Let() (ast.Let, diag.Error) {
	p.advance() // skip "let"

	if p.c.Kind != token.Word {
		return ast.Let{}, p.unexpected()
	}
	if p.n.Kind == token.Walrus {
		panic("not implemented")
		// return p.letWalrusStatement()
	}
	name := p.word()

	if p.c.Kind != token.Colon {
		return ast.Let{}, p.unexpected()
	}
	p.advance() // skip ":"

	typ, err := p.TypeSpec()
	if err != nil {
		return ast.Let{}, err
	}

	if p.c.Kind != token.Assign {
		return ast.Let{}, p.unexpected()
	}
	p.advance() // skip "="

	exp, err := p.Exp()
	if err != nil {
		return ast.Let{}, err
	}

	if p.c.Kind != token.Semicolon {
		return ast.Let{}, p.unexpected()
	}
	p.advance() // consume ";"

	return ast.Let{
		Name: name,
		Type: typ,
		Exp:  exp,
	}, nil
}
