package parser

import (
	"github.com/mebyus/ku/internal/ku/ast"
	"github.com/mebyus/ku/internal/ku/token"
)

func (p *Parser) assignSymbol() (ast.Statement, ss) {
	name := p.peek.Data
	pin := p.peek.Pin

	p.advance() // skip symbol name
	p.advance() // skip "="

	a := ast.AssignSymbol{
		Pin:  pin,
		Name: name,
	}

	exp, s := p.Exp()
	a.Exp = exp
	if s != 0 {
		return &a, s
	}

	if p.peek.Kind != token.Semicolon {
		p.report(p.peek.Pin, "missing \";\" after assign statement")
	} else {
		p.advance() // skip ";"
	}

	return &a, 0
}
