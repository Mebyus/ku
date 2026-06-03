package parser

import (
	"fmt"

	"github.com/mebyus/ku/internal/ku/ast"
	"github.com/mebyus/ku/internal/ku/enums/bop"
	"github.com/mebyus/ku/internal/ku/enums/uop"
	"github.com/mebyus/ku/internal/ku/token"
)

func (p *Parser) Exp() ast.Exp {
	return p.pratt(0)
}

func (p *Parser) pratt(power int) ast.Exp {
	a := p.primary()
	if a == nil {
		// TODO: should be an error here?
		return nil
	}

	for {
		k, ok := bop.FromToken(p.peek.Kind)
		if !ok || k.Power() <= power {
			return a
		}
		op := bop.Op{Pin: p.peek.Pin, Kind: k}
		p.advance() // skip binary operator

		b := p.pratt(k.Power())
		if b == nil {
			// TODO: should be an error here?
			return nil
		}

		a = &ast.BinExp{Op: op, A: a, B: b}
	}
}

func (p *Parser) primary() ast.Exp {
	k, ok := uop.FromToken(p.peek.Kind)
	if !ok {
		return p.operand()
	}

	op := uop.Op{Pin: p.peek.Pin, Kind: k}
	p.advance() // skip unary operator

	exp := p.primary()
	if exp == nil {
		return nil
	}
	return &ast.UnExp{
		Op: op,
		A:  exp,
	}
}

func (p *Parser) operand() ast.Operand {
	switch p.peek.Kind {
	case token.Integer:
		tok := p.peek
		p.advance() // skip integer
		return &ast.Integer{
			Pin: tok.Pin,
			Val: tok.Val,
			Aux: tok.Flags,
		}
	default:
		pin := p.peek.Pin
		er := ast.Error{
			Pin:   pin,
			Short: fmt.Sprintf("expected expression or operand, found %s token instead", p.peek.Kind),
		}
		p.addError(&er)
		p.advance() // TODO: we should do a sync here
		return &ast.ErrorExp{Error: er}
	}
}
