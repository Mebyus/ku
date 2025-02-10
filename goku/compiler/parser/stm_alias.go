package parser

import (
	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/token"
)

func (p *Parser) topAlias(traits ast.Traits) diag.Error {
	a, err := p.Alias()
	if err != nil {
		return err
	}

	p.text.AddAlias(ast.TopAlias{
		Alias:  a,
		Traits: traits,
	})
	return nil
}

func (p *Parser) Alias() (ast.Alias, diag.Error) {
	p.advance() // skip "let"

	if p.c.Kind != token.Word {
		return ast.Alias{}, p.unexpected()
	}
	name := p.word()

	if p.c.Kind != token.RightArrow {
		return ast.Alias{}, p.unexpected()
	}
	p.advance() // skip "=>"

	exp, err := p.Exp()
	if err != nil {
		return ast.Alias{}, err
	}

	if p.c.Kind != token.Semicolon {
		return ast.Alias{}, p.unexpected()
	}
	p.advance() // consume ";"

	return ast.Alias{
		Name: name,
		Exp:  exp,
	}, nil
}
