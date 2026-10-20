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

// parse statement which starts as
//
//	word ".*"
func (p *Parser) symbolDerefStart() (ast.Statement, ss) {
	name := p.peek.Data
	pin := p.peek.Pin

	p.advance() // skip symbol name
	p.advance() // skip ".*"

	if p.peek.Kind != token.Assign {
		panic("stub")
	}

	p.advance() // skip "="

	exp, s := p.Exp()
	if s != 0 {
		return nil, s
	}

	if p.peek.Kind != token.Semicolon {
		p.report(p.peek.Pin, "missing \";\" after assign statement")
	} else {
		p.advance() // skip ";"
	}

	return &ast.Assign{
		Target: &ast.SymDeref{
			Name: name,
			Pin:  pin,
		},
		Exp: exp,
	}, 0
}
