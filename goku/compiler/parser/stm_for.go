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

func (p *Parser) forRange() (ast.Statement, diag.Error) {
	name := p.word()
	p.advance() // skip ":"

	typ, err := p.TypeSpec()
	if err != nil {
		return nil, err
	}

	if p.c.Kind != token.In {
		return nil, p.unexpected()
	}
	p.advance() // skip "in"

	if p.c.Kind != token.Word && p.c.Data != "range" {
		return nil, p.unexpected()
	}
	p.advance() // skip "range"

	if p.c.Kind != token.LeftParen {
		return nil, p.unexpected()
	}
	p.advance() // skip "("

	exp, err := p.Exp()
	if err != nil {
		return nil, err
	}

	if p.c.Kind == token.Comma {
		return p.continueForRange2(name, typ, exp)
	}

	if p.c.Kind != token.RightParen {
		return nil, p.unexpected()
	}
	p.advance() // skip ")"

	body, err := p.Block()
	if err != nil {
		return nil, err
	}

	return ast.ForRange{
		Name: name,
		Exp:  exp,
		Type: typ,
		Body: body,
	}, nil
}

func (p *Parser) continueForRange2(name ast.Word, spec ast.TypeSpec, start ast.Exp) (ast.ForRange2, diag.Error) {
	p.advance() // skip ","

	end, err := p.Exp()
	if err != nil {
		return ast.ForRange2{}, err
	}

	if p.c.Kind != token.RightParen {
		return ast.ForRange2{}, p.unexpected()
	}
	p.advance() // skip ")"

	body, err := p.Block()
	if err != nil {
		return ast.ForRange2{}, err
	}

	return ast.ForRange2{
		Name:  name,
		Type:  spec,
		Start: start,
		End:   end,
		Body:  body,
	}, nil
}
