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

	if p.peek.Kind != token.Word {
		return ast.Const{}, p.unexpected()
	}
	if p.next.Kind == token.Walrus {
		return p.walrusConst()
	}
	name := p.word()

	if p.peek.Kind != token.Colon {
		return ast.Const{}, p.unexpected()
	}
	p.advance() // skip ":"

	typ, err := p.TypeSpec()
	if err != nil {
		return ast.Const{}, err
	}

	if p.peek.Kind != token.Assign {
		return ast.Const{}, p.unexpected()
	}
	p.advance() // skip "="

	exp, err := p.Exp()
	if err != nil {
		return ast.Const{}, err
	}

	if p.peek.Kind != token.Semicolon {
		return ast.Const{}, p.unexpected()
	}
	p.advance() // consume ";"

	return ast.Const{
		Name: name,
		Type: typ,
		Exp:  exp,
	}, nil
}

func (p *Parser) walrusConst() (ast.Const, diag.Error) {
	name := p.word()
	p.advance() // skip ":="

	exp, err := p.Exp()
	if err != nil {
		return ast.Const{}, err
	}

	if p.peek.Kind != token.Semicolon {
		return ast.Const{}, p.unexpected()
	}
	p.advance() // consume ";"

	return ast.Const{
		Name: name,
		Exp:  exp,
	}, nil
}
