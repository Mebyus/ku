package parser

import (
	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/token"
)

func (p *Parser) Statement() (ast.Statement, diag.Error) {
	switch p.c.Kind {
	case token.LeftCurly:
		return p.Block()
	case token.Let:
		return p.Let()
	case token.Var:
		return p.Var()
	case token.If:
		return p.If()
	case token.Ret:
		return p.Ret()
	case token.For:
		return p.For()
	// case token.Jump:
	// 	return p.jumpStatement()
	// case token.Never:
	// 	return p.neverStatement()
	case token.Stub:
		return p.Stub()
	// case token.Defer:
	// 	return p.deferStatement()
	case token.Word:
		return p.AssignOrInvoke()
	default:
		return nil, p.unexpected()
	}
}

func (p *Parser) Block() (ast.Block, diag.Error) {
	pin := p.c.Pin
	p.advance() // skip "{"

	var nodes []ast.Statement
	for {
		if p.c.Kind == token.RightCurly {
			p.advance() // skip "}"
			return ast.Block{
				Pin:   pin,
				Nodes: nodes,
			}, nil
		}

		s, err := p.Statement()
		if err != nil {
			return ast.Block{}, err
		}
		nodes = append(nodes, s)
	}
}
