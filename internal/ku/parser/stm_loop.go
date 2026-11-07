package parser

import (
	"fmt"

	"github.com/mebyus/ku/internal/ku/ast"
	"github.com/mebyus/ku/internal/ku/token"
)

func (p *Parser) For() (ast.Statement, ss) {
	if p.next.Kind == token.LeftCurly {
		return p.loop()
	}
	if p.next.Kind == token.LeftSquare {
		return p.fornum()
	}

	pin := p.peek.Pin
	p.advance() // skip "for"

	exp, s := p.Exp()
	if s != 0 {
		return nil, s
	}

	f := ast.While{Exp: exp, Pin: pin}

	if p.peek.Kind != token.LeftCurly {
		p.report(p.peek.Pin, 0, fmt.Sprintf("expected \"{\" to start loop body in \"for\" statement, found %s token", &p.peek))
		return &f, ssNode
	}
	s = p.block(&f.Body)
	if s != 0 {
		return &f, s
	}
	return &f, 0
}

func (p *Parser) loop() (ast.Statement, ss) {
	pin := p.peek.Pin
	p.advance() // skip "for"

	loop := ast.Loop{Pin: pin}

	if p.peek.Kind != token.LeftCurly {
		p.report(p.peek.Pin, 0, fmt.Sprintf("expected \"{\" to start loop body in \"for\" statement, found %s token", &p.peek))
		return &loop, ssNode
	}
	s := p.block(&loop.Body)
	if s != 0 {
		return &loop, s
	}
	return &loop, 0
}

func (p *Parser) fornum() (ast.Statement, ss) {
	pin := p.peek.Pin
	p.advance() // skip "for"

	p.advance() // skip "["

	f := ast.ForNum{Pin: pin}

	if p.peek.Kind != token.Colon {
		a, s := p.Exp()
		if s != 0 {
			return nil, s
		}
		f.A = a
	}
	if p.peek.Kind == token.Colon {
		p.advance() // skip ":"
	} else {
		p.report(p.peek.Pin, 0, fmt.Sprintf("expected \":\" in for range, found %s token instead", &p.peek))
	}

	if p.peek.Kind == token.RightSquare {
		p.report(p.peek.Pin, 0, "missing end expression in for range statement")
		f.B = &ast.InvExp{Pin: p.peek.Pin}
		p.advance() // skip "]"
	} else {
		b, s := p.Exp()
		if s != 0 {
			return nil, s
		}
		f.B = b
		if p.peek.Kind != token.RightSquare {
			p.report(p.peek.Pin, 0, "missing \"]\" after end expression in for range")
		} else {
			p.advance() // skip "]"
		}
	}

	if p.peek.Kind != token.LeftCurly {
		p.report(p.peek.Pin, 0, fmt.Sprintf("expected \"{\" to start loop body in \"for range\" statement, found %s token", &p.peek))
		return &f, ssNode
	}
	s := p.block(&f.Body)
	if s != 0 {
		return &f, s
	}
	return &f, 0
}
