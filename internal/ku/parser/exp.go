package parser

import (
	"fmt"

	"github.com/mebyus/ku/internal/ku/ast"
	"github.com/mebyus/ku/internal/ku/enums/bop"
	"github.com/mebyus/ku/internal/ku/enums/uop"
	"github.com/mebyus/ku/internal/ku/token"
)

func (p *Parser) Exp() (ast.Exp, ss) {
	return p.pratt(0)
}

func (p *Parser) pratt(power int) (ast.Exp, ss) {
	a, s := p.primary()
	if s != 0 {
		return a, s
	}

	for {
		k, ok := bop.FromToken(p.peek.Kind)
		if !ok || k.Power() <= power {
			return a, 0
		}

		op := bop.Op{Pin: p.peek.Pin, Kind: k}
		p.advance() // skip binary operator

		b, s := p.pratt(k.Power())
		if s != 0 {
			return a, s
		}

		a = &ast.BinExp{Op: op, A: a, B: b}
	}
}

func (p *Parser) primary() (ast.Exp, ss) {
	k, ok := uop.FromToken(p.peek.Kind)
	if !ok {
		return p.operand()
	}

	op := uop.Op{Pin: p.peek.Pin, Kind: k}
	p.advance() // skip unary operator

	exp, s := p.primary()
	return &ast.UnExp{
		Op: op,
		A:  exp,
	}, s
}

func (p *Parser) operand() (ast.Operand, ss) {
	switch p.peek.Kind {
	case token.Integer:
		tok := p.peek
		p.advance() // skip integer
		return &ast.Integer{
			Pin: tok.Pin,
			Val: tok.Val,
			Aux: tok.Flags,
		}, 0
	case token.True:
		tok := p.peek
		p.advance() // skip "true"
		return &ast.True{Pin: tok.Pin}, 0
	case token.False:
		tok := p.peek
		p.advance() // skip "false"
		return &ast.False{Pin: tok.Pin}, 0
	case token.Word:
		return p.Chain()
	case token.LeftParen:
		return p.parenExp()
	}

	return p.syncExp(p.peek.Pin, fmt.Sprintf("expected expression or operand, found %s token instead", p.peek.Kind))
}

func (p *Parser) parenExp() (ast.Operand, ss) {
	pin := p.peek.Pin
	p.advance() // skip "("

	exp, s := p.Exp()
	if s != 0 {
		return &ast.ParenExp{
			Pin: pin,
			Exp: exp,
		}, s
	}
	if p.peek.Kind == token.RightParen {
		p.advance() // skip ")"
		return &ast.ParenExp{
			Pin: pin,
			Exp: exp,
		}, 0
	}

	pin = p.peek.Pin
	er := ast.Error{
		Pin:   pin,
		Short: fmt.Sprintf("expected \")\" to close parenthesized expression, found %s token instead", p.peek.Kind),
	}
	p.addError(&er)
	p.advance() // TODO: we should do a sync here
	return &ast.ErrorExp{Error: er}, 0
}
