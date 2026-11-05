package parser

import (
	"github.com/mebyus/ku/internal/ku/ast"
	"github.com/mebyus/ku/internal/ku/token"
)

func (p *Parser) invokeSymbol() (ast.Statement, ss) {
	name := p.peek.Data
	pin := p.peek.Pin

	p.advance() // skip symbol name

	args, s := p.args()
	if s != 0 {
		return nil, s
	}

	switch p.peek.Kind {
	case token.Assign:
		s = ssNode
		p.report(p.peek.Pin, 0, "cannot assign to call expression")
	case token.Semicolon:
		p.advance() // skip ";"
	default:
		p.report(p.peek.Pin, 0, "missing \";\" after invoke")
	}

	return &ast.InvokeSymbol{
		Args: args,
		Name: name,
		Pin:  pin,
	}, s
}
