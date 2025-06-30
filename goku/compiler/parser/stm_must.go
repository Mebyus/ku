package parser

import (
	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/token"
)

func (p *Parser) StaticMust() (ast.StaticMust, diag.Error) {
	p.advance() // skip "#must"

	if p.c.Kind != token.LeftParen {
		return ast.StaticMust{}, p.unexpected()
	}
	p.advance() // skip "("

	exp, err := p.Exp()
	if err != nil {
		return ast.StaticMust{}, err
	}

	if p.c.Kind != token.RightParen {
		return ast.StaticMust{}, p.unexpected()
	}
	p.advance() // skip ")"

	if p.c.Kind != token.Semicolon {
		return ast.StaticMust{}, p.unexpected()
	}
	p.advance() // skip ";"

	return ast.StaticMust{Exp: exp}, nil
}

func (p *Parser) topStaticMust() diag.Error {
	must, err := p.StaticMust()
	if err != nil {
		return err
	}
	p.text.AddMust(must)
	return nil
}
