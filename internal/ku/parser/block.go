package parser

import (
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
	default:
		panic("stub")
	}
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
