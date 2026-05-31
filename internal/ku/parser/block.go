package parser

import (
	"fmt"

	"github.com/mebyus/ku/internal/ku/ast"
	"github.com/mebyus/ku/internal/ku/token"
)

func (p *Parser) Block() *ast.Block {
	if p.peek.Kind != token.LeftCurly {
		p.report(p.peek.Pin, fmt.Sprintf("expected \"{\" as block start, found %s token instead", p.peek.Kind))
		return nil
	}

	pin := p.peek.Pin
	p.advance() // skip "{"

	var nodes []ast.Statement
	for {
		if p.peek.Kind == token.RightCurly {
			p.advance() // skip "}"
			return &ast.Block{
				Pin:   pin,
				Nodes: nodes,
			}
		}

		node, ok := p.Statement()
		if !ok {
			return nil
		}
		if node != nil {
			nodes = append(nodes, node)
		}
	}
}

func (p *Parser) Statement() (ast.Statement, bool) {
	switch p.peek.Kind {
	case token.Return:
		return p.Return()
	default:
		panic("stub")
	}
}

func (p *Parser) Return() (ast.Statement, bool) {
	pin := p.peek.Pin
	p.advance() // skip "return"

	if p.peek.Kind == token.Semicolon {
		p.advance() // skip ";"
		return &ast.Return{Pin: pin}, true
	}

	exp := p.Exp()

	if p.peek.Kind != token.Semicolon {
		p.report(p.peek.Pin, "missing \";\" after return statement")
		// continue parsing, not a serious error
	} else {
		p.advance() // skip ";"
	}

	return &ast.Return{
		Pin: pin,
		Exp: exp,
	}, true
}
