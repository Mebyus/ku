package parser

import (
	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/token"
)

func (p *Parser) Statement() (ast.Statement, diag.Error) {
	switch p.peek.Kind {
	case token.LeftCurly:
		return p.Block()
	case token.Const:
		return p.Const()
	case token.Var:
		return p.Var()
	case token.Let:
		return p.Alias()
	case token.If:
		return p.If()
	case token.Ret:
		return p.Ret()
	case token.For:
		return p.For()
	case token.Jump:
		return p.Jump()
	case token.Never:
		return p.Never()
	case token.Stub:
		return p.Stub()
	case token.Panic:
		return p.Panic()
	case token.Must:
		return p.Must()
	case token.StaticMust:
		return p.StaticMust()
	case token.StaticIf:
		return p.StaticIf()
	case token.Test:
		return p.Test()
	case token.Debug:
		return p.Debug()
	// case token.Defer:
	// 	return p.deferStatement()
	case token.Word, token.Unsafe, token.Underscore:
		return p.AssignOrInvoke()
	default:
		return nil, p.unexpected()
	}
}

func (p *Parser) Block() (ast.Block, diag.Error) {
	if p.peek.Kind != token.LeftCurly {
		return ast.Block{}, p.unexpected()
	}
	pin := p.peek.Pin
	p.advance() // skip "{"

	var nodes []ast.Statement
	for {
		if p.peek.Kind == token.RightCurly {
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

func (p *Parser) Static() (ast.Static, diag.Error) {
	if p.peek.Kind != token.HashCurly {
		return ast.Static{}, p.unexpected()
	}
	pin := p.peek.Pin
	p.advance() // skip "#{"

	var nodes []ast.Statement
	for {
		if p.peek.Kind == token.RightCurly {
			p.advance() // skip "}"
			return ast.Static{
				Pin:   pin,
				Nodes: nodes,
			}, nil
		}

		s, err := p.Statement()
		if err != nil {
			return ast.Static{}, err
		}
		nodes = append(nodes, s)
	}
}

func (p *Parser) Debug() (ast.Debug, diag.Error) {
	p.advance() // skip "#debug"

	if p.peek.Kind != token.LeftCurly {
		return ast.Debug{}, p.unexpected()
	}

	block, err := p.Block()
	if err != nil {
		return ast.Debug{}, err
	}

	return ast.Debug{Block: block}, nil
}
