package parser

import (
	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/token"
)

func (p *Parser) topConst(traits ast.Traits) diag.Error {
	l, err := p.Const()
	if err != nil {
		return err
	}

	p.text.AddConst(ast.TopConst{
		Const:  l,
		Traits: traits,
	})
	return nil
}

func (p *Parser) Const() (ast.Const, diag.Error) {
	p.advance() // skip "const"

	if p.c.Kind != token.Word {
		return ast.Const{}, p.unexpected()
	}
	if p.n.Kind == token.Walrus {
		panic("not implemented")
		// return p.letWalrusStatement()
	}
	name := p.word()

	if p.c.Kind != token.Colon {
		return ast.Const{}, p.unexpected()
	}
	p.advance() // skip ":"

	typ, err := p.TypeSpec()
	if err != nil {
		return ast.Const{}, err
	}

	if p.c.Kind != token.Assign {
		return ast.Const{}, p.unexpected()
	}
	p.advance() // skip "="

	exp, err := p.Exp()
	if err != nil {
		return ast.Const{}, err
	}

	if p.c.Kind != token.Semicolon {
		return ast.Const{}, p.unexpected()
	}
	p.advance() // consume ";"

	return ast.Const{
		Name: name,
		Type: typ,
		Exp:  exp,
	}, nil
}
