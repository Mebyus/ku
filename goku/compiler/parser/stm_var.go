package parser

import (
	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/token"
)

func (p *Parser) topVar(traits ast.Traits) diag.Error {
	v, err := p.Var()
	if err != nil {
		return err
	}

	p.text.AddVar(ast.TopVar{
		Var:    v,
		Traits: traits,
	})
	return nil
}

func (p *Parser) Var() (ast.Var, diag.Error) {
	p.advance() // skip "var"
	if p.c.Kind != token.Word {
		return ast.Var{}, p.unexpected()
	}
	name := p.word()

	if p.c.Kind != token.Colon {
		return ast.Var{}, p.unexpected()
	}
	p.advance() // skip ":"

	typ, err := p.TypeSpec()
	if err != nil {
		return ast.Var{}, err
	}

	if p.c.Kind == token.Semicolon {
		// empty init expression

		p.advance() // skip ";"
		return ast.Var{
			Name: name,
			Type: typ,
		}, nil
	}

	if p.c.Kind != token.Assign {
		return ast.Var{}, p.unexpected()
	}
	p.advance() // skip "="

	exp, err := p.InitExp()
	if err != nil {
		return ast.Var{}, err
	}

	if p.c.Kind != token.Semicolon {
		return ast.Var{}, p.unexpected()
	}
	p.advance() // consume ";"

	return ast.Var{
		Name: name,
		Type: typ,
		Exp:  exp,
	}, nil
}
