package parser

import (
	"fmt"

	"github.com/mebyus/ku/internal/ku/ast"
	"github.com/mebyus/ku/internal/ku/token"
)

func (p *Parser) Goto() (ast.Statement, ss) {
	p.advance() // skip "goto"

	pin := p.peek.Pin
	if p.peek.Kind != token.Word {
		p.report(pin, 0, fmt.Sprintf("expected label name in goto statement, found %s token instead", &p.peek))
		// TODO: sync
		return nil, ssNode
	}
	name := p.peek.Data
	p.advance() // skip label name

	if p.peek.Kind != token.Semicolon {
		p.report(p.peek.Pin, 0, "missing \";\" after goto statement")
	} else {
		p.advance() // skip ";"
	}

	return &ast.Goto{
		Name: name,
		Pin:  pin,
	}, 0
}

func (p *Parser) Zone() (ast.Statement, ss) {
	p.advance() // skip "zone"

	pin := p.peek.Pin
	if p.peek.Kind != token.Word {
		p.report(pin, 0, fmt.Sprintf("expected zone name, found %s token instead", &p.peek))
		// TODO: sync
		return nil, ssNode
	}
	name := p.peek.Data
	p.advance() // skip zone name

	if p.peek.Kind == token.LeftCurly {
		var block ast.Block
		s := p.block(&block)
		return &ast.ZoneBlock{
			Nodes: block.Nodes,
			Name:  name,
			Pin:   pin,
		}, s
	}

	if p.peek.Kind != token.Semicolon {
		p.report(p.peek.Pin, 0, "missing \";\" after place statement")
	} else {
		p.advance() // skip ";"
	}

	return &ast.Place{
		Name: name,
		Pin:  pin,
	}, 0
}
