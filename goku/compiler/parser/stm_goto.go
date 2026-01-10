package parser

import (
	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/token"
)

func (p *Parser) Goto() (ast.Goto, diag.Error) {
	p.advance() // skip "goto"

	if p.peek.Kind != token.Label {
		return ast.Goto{}, p.unexpected()
	}
	pin := p.peek.Pin
	name := p.peek.Data

	if p.peek.Kind != token.Semicolon {
		return ast.Goto{}, p.unexpected()
	}
	p.advance() // skip ";"

	return ast.Goto{
		Name: name,
		Pin:  pin,
	}, nil
}

func (p *Parser) Gonext() (ast.Gonext, diag.Error) {
	pin := p.peek.Pin
	p.advance() // skip "gonext"

	if p.peek.Kind != token.Semicolon {
		return ast.Gonext{}, p.unexpected()
	}
	p.advance() // skip ";"

	return ast.Gonext{Pin: pin}, nil
}

func (p *Parser) Break() (ast.Break, diag.Error) {
	pin := p.peek.Pin
	p.advance() // skip "break"

	if p.peek.Kind != token.Semicolon {
		return ast.Break{}, p.unexpected()
	}
	p.advance() // skip ";"

	return ast.Break{Pin: pin}, nil
}
