package parser

import (
	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/enums/bok"
	"github.com/mebyus/ku/goku/compiler/enums/uok"
)

// Parse single expression (no expression will result in error).
// Single means pack expression is not recognized.
//
// Parsing is done via Pratt's recursive descent algorithm variant.
func (p *Parser) Exp() (ast.Exp, diag.Error) {
	return p.pratt(0)
}

func (p *Parser) pratt(power int) (ast.Exp, diag.Error) {
	a, err := p.Primary()
	if err != nil {
		return nil, err
	}

	for {
		k, ok := bok.FromToken(p.c.Kind)
		if !ok || k.Power() <= power {
			return a, nil
		}
		op := ast.BinOp{Pin: p.c.Pin, Kind: k}
		p.advance() // skip binary operator

		b, err := p.pratt(k.Power())
		if err != nil {
			return nil, err
		}

		a = ast.Binary{Op: op, A: a, B: b}
	}
}

func (p *Parser) Primary() (ast.Exp, diag.Error) {
	k, ok := uok.FromToken(p.c.Kind)
	if !ok {
		return p.Operand()
	}

	op := ast.UnaryOp{Pin: p.c.Pin, Kind: k}
	p.advance() // skip unary operator

	exp, err := p.Primary()
	if err != nil {
		return nil, err
	}
	return ast.Unary{
		Op:  op,
		Exp: exp,
	}, nil
}
