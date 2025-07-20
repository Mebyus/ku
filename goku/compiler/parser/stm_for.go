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
	if p.c.Kind == token.Word {
		switch p.n.Kind {
		case token.Assign:
			return p.forRangeAutoType()
		case token.Colon:
			return p.forRange()
		}
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

func (p *Parser) forRangeAutoType() (ast.Statement, diag.Error) {
	name := p.word()
	p.advance() // skip "="

	start, end, err := p.forRangeStartEnd()
	if err != nil {
		return nil, err
	}

	body, err := p.Block()
	if err != nil {
		return nil, err
	}

	return ast.ForRange{
		Name:  name,
		Start: start,
		End:   end,
		Body:  body,
	}, nil
}

func (p *Parser) forRange() (ast.Statement, diag.Error) {
	name := p.word()
	p.advance() // skip ":"

	typ, err := p.TypeSpec()
	if err != nil {
		return nil, err
	}

	if p.c.Kind != token.Assign {
		return nil, p.unexpected()
	}
	p.advance() // skip "="

	start, end, err := p.forRangeStartEnd()
	if err != nil {
		return nil, p.unexpected()
	}

	body, err := p.Block()
	if err != nil {
		return nil, err
	}

	return ast.ForRange{
		Name:  name,
		Start: start,
		End:   end,
		Type:  typ,
		Body:  body,
	}, nil
}

func (p *Parser) forRangeStartEnd() (ast.Exp, ast.Exp, diag.Error) {
	if p.c.Kind != token.LeftSquare {
		return nil, nil, p.unexpected()
	}
	p.advance() // skip "["

	var start ast.Exp
	if p.c.Kind != token.Colon {
		exp, err := p.Exp()
		if err != nil {
			return nil, nil, err
		}
		start = exp
	}

	if p.c.Kind != token.Colon {
		return nil, nil, p.unexpected()
	}
	p.advance() // skip ":"

	end, err := p.Exp()
	if err != nil {
		return nil, nil, err
	}

	if p.c.Kind != token.RightSquare {
		return nil, nil, p.unexpected()
	}
	p.advance() // skip "["

	return start, end, nil
}
