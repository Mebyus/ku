package parser

import (
	"fmt"

	"github.com/mebyus/ku/internal/ku/ast"
	"github.com/mebyus/ku/internal/ku/token"
)

func (p *Parser) block(b *ast.Block) ss {
	b.Pin = p.peek.Pin
	p.advance() // skip "{"

	for {
		if p.peek.Kind == token.RightCurly {
			p.advance() // skip "}"
			return 0
		}

		node, s := p.Statement()
		if node != nil {
			b.Nodes = append(b.Nodes, node)
		}
		switch s {
		case ssTop, ssAbort:
			return s
		}
	}
}

func (p *Parser) Statement() (ast.Statement, ss) {
	switch p.peek.Kind {
	case token.Return:
		return p.Return()
	case token.Const:
		return p.Const()
	default:
		panic("stub")
	}
}

func (p *Parser) Const() (ast.Statement, ss) {
	p.advance() // skip "const"

	if p.peek.Kind != token.Word {
		p.report(p.peek.Pin, fmt.Sprintf("expected constant name, found %s token instead", p.peek.Kind))
		// TODO: sync
		return nil, ssNode
	}

	name := p.peek.Data
	pin := p.peek.Pin
	p.advance() // skip constant name

	var typ ast.TypeSpec
	switch p.peek.Kind {
	case token.Walrus:
		// untyped constant
		p.advance() // skip ":="
	case token.Colon:
		p.advance() // skip ":"
		if p.peek.Kind == token.Assign {
			// ":" + "=" same as ":="
			p.advance() // skip "="
		} else {
			t, s := p.TypeSpec()
			if s != 0 {
				return &ast.Const{
					Pin:  pin,
					Name: name,
					Type: t,
				}, s
			}
			typ = t

			if p.peek.Kind != token.Assign {
				p.report(p.peek.Pin, "missing \"=\" before constant value expression")
			} else {
				p.advance() // skip "="
			}
		}
	case token.Assign:
		p.report(p.peek.Pin, "assign \"=\" operator used instead of \":=\" for constant definition")
		p.advance() // skip "="
	default:
		p.report(p.peek.Pin, fmt.Sprintf("expected \":\" or \":=\" after constant name, found %s token instead", p.peek.Kind))
		// TODO: sync?
		p.advance()
	}

	var exp ast.Exp
	switch p.peek.Kind {
	case token.Semicolon, token.RightCurly:
		p.report(p.peek.Pin, "missing constant value expression")
	default:
		e, s := p.Exp()
		if s != 0 {
			return &ast.Const{
				Pin:  pin,
				Name: name,
				Type: typ,
				Exp:  e,
			}, s
		}
		exp = e
	}

	if p.peek.Kind != token.Semicolon {
		p.report(p.peek.Pin, "missing \";\" after constant definition")
	} else {
		p.advance() // skip ";"
	}

	return &ast.Const{
		Pin:  pin,
		Name: name,
		Type: typ,
		Exp:  exp,
	}, 0
}

func (p *Parser) Return() (ast.Statement, ss) {
	pin := p.peek.Pin
	p.advance() // skip "return"

	switch p.peek.Kind {
	case token.Semicolon:
		p.advance() // skip ";"
		return &ast.Return{Pin: pin}, 0
	case token.RightCurly:
		p.report(pin, "missing \";\" after return statement")
		return &ast.Return{Pin: pin}, 0
	}

	exp, s := p.Exp()
	if s != 0 {
		return &ast.Return{
			Pin: pin,
			Exp: exp,
		}, s
	}

	if p.peek.Kind != token.Semicolon {
		p.report(p.peek.Pin, "missing \";\" after return statement")
		// continue parsing, not a serious error
	} else {
		p.advance() // skip ";"
	}

	return &ast.Return{
		Pin: pin,
		Exp: exp,
	}, 0
}
