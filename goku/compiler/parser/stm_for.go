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

	return p.While()
}

func (p *Parser) Loop() (ast.Loop, diag.Error) {
	p.advance() // skip "for"

	body, err := p.Block()
	if err != nil {
		return ast.Loop{}, err
	}

	return ast.Loop{Body: body}, nil
}

func (p *Parser) While() (ast.While, diag.Error) {
	p.advance() // skip "for"

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
