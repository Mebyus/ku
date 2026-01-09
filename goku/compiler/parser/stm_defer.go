package parser

import (
	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/token"
)

func (p *Parser) DeferCall() (ast.DeferCall, diag.Error) {
	p.advance() // skip "defer"

	o, err := p.Chain()
	if err != nil {
		return ast.DeferCall{}, err
	}

	c, ok := o.(ast.Call)
	if !ok {
		return ast.DeferCall{}, &diag.SimpleMessageError{
			Pin:  o.Span().Pin,
			Text: "defer only accepts block statement or call expression",
		}
	}

	if p.peek.Kind != token.Semicolon {
		return ast.DeferCall{}, p.unexpected()
	}
	p.advance() // skip ";"

	return ast.DeferCall{Call: c}, nil
}

func (p *Parser) DeferBlock() (ast.DeferCall, diag.Error) {
	p.advance() // skip "defer"

	block, err := p.Block()
	if err != nil {
		return ast.DeferCall{}, err
	}

	_ = block

	panic("not implemented")
}
