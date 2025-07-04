package parser

import (
	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/token"
)

func (p *Parser) For() (ast.Statement, diag.Error) {
	if p.n.Kind == token.LeftCurly {
		return p.Loop()
	}

	p.advance() // skip "for"
	if p.c.Kind == token.Word && p.n.Kind == token.Colon {
		return p.forRange()
	}

	return p.while()
}

func (p *Parser) Loop() (ast.Loop, diag.Error) {
	p.advance() // skip "for"

	body, err := p.Block()
	if err != nil {
		return ast.Loop{}, err
	}

	return ast.Loop{Body: body}, nil
}

func (p *Parser) while() (ast.While, diag.Error) {
	exp, err := p.Exp()
	if err != nil {
		return ast.While{}, err
	}

	body, err := p.Block()
	if err != nil {
		return ast.While{}, err
	}

	return ast.While{
		Body: body,
		Exp:  exp,
	}, nil
}

func (p *Parser) forRange() (ast.ForRange, diag.Error) {
	name := p.word()
	p.advance() // skip ":"

	typ, err := p.TypeSpec()
	if err != nil {
		return ast.ForRange{}, err
	}

	if p.c.Kind != token.In {
		return ast.ForRange{}, p.unexpected()
	}
	p.advance() // skip "in"

	if p.c.Kind != token.Word && p.c.Data != "range" {
		return ast.ForRange{}, p.unexpected()
	}
	p.advance() // skip "range"

	if p.c.Kind != token.LeftParen {
		return ast.ForRange{}, p.unexpected()
	}
	p.advance() // skip "("

	exp, err := p.Exp()
	if err != nil {
		return ast.ForRange{}, err
	}

	if p.c.Kind != token.RightParen {
		return ast.ForRange{}, p.unexpected()
	}
	p.advance() // skip ")"

	body, err := p.Block()
	if err != nil {
		return ast.ForRange{}, err
	}

	return ast.ForRange{
		Name: name,
		Exp:  exp,
		Type: typ,
		Body: body,
	}, nil
}
