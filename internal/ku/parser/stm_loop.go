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

	if p.peek.Kind == token.Word && p.next.Kind == token.Walrus {
		return p.forrange()
	}

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

	f := ast.ForNum{Pin: pin}

	a, b, s := p.loopRange()
	f.A = a
	f.B = b
	if s != 0 {
		return nil, s
	}

	if p.peek.Kind != token.LeftCurly {
		p.report(p.peek.Pin, 0, fmt.Sprintf("expected \"{\" to start loop body in \"for range\" statement, found %s token", &p.peek))
		return &f, ssNode
	}
	s = p.block(&f.Body)
	if s != 0 {
		return &f, s
	}
	return &f, 0
}

func (p *Parser) forrange() (ast.Statement, ss) {
	pin := p.peek.Pin
	name := p.peek.Data
	p.advance() // skip loop variable name

	p.advance() // skip ":="

	f := ast.ForRange{
		Pin:  pin,
		Name: name,
	}

	a, b, s := p.loopRange()
	f.A = a
	f.B = b
	if s != 0 {
		return nil, s
	}

	if p.peek.Kind != token.LeftCurly {
		p.report(p.peek.Pin, 0, fmt.Sprintf("expected \"{\" to start loop body in \"for range\" statement, found %s token", &p.peek))
		return &f, ssNode
	}
	s = p.block(&f.Body)
	if s != 0 {
		return &f, s
	}
	return &f, 0
}

func (p *Parser) loopRange() (ast.Exp, ast.Exp, ss) {
	if p.peek.Kind != token.LeftSquare {
		p.report(p.peek.Pin, 0, fmt.Sprintf("expected \"[\" in for range expression, found %s token instead", &p.peek))
		return nil, nil, ssNode
	}
	p.advance() // skip "["

	var a, b ast.Exp
	if p.peek.Kind != token.Colon {
		e, s := p.Exp()
		if s != 0 {
			return nil, nil, s
		}
		a = e
	}
	if p.peek.Kind == token.Colon {
		p.advance() // skip ":"
	} else {
		p.report(p.peek.Pin, 0, fmt.Sprintf("expected \":\" in for range, found %s token instead", &p.peek))
	}

	if p.peek.Kind == token.RightSquare {
		p.report(p.peek.Pin, 0, "missing end expression in for range statement")
		b = &ast.InvExp{Pin: p.peek.Pin}
		p.advance() // skip "]"
	} else {
		e, s := p.Exp()
		if s != 0 {
			return a, nil, s // TODO: invalid expression for B?
		}
		b = e
		if p.peek.Kind != token.RightSquare {
			p.report(p.peek.Pin, 0, "missing \"]\" after end expression in for range")
		} else {
			p.advance() // skip "]"
		}
	}

	return a, b, 0
}
