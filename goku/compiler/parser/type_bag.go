package parser

import (
	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/token"
)

func (p *Parser) Bag() (ast.Bag, diag.Error) {
	pin := p.peek.Pin
	p.advance() // skip "bag"

	if p.peek.Kind != token.LeftCurly {
		return ast.Bag{}, p.unexpected()
	}
	p.advance() // skip "{"

	var funs []ast.BagFun
	for {
		if p.peek.Kind == token.RightCurly {
			p.advance() // skip "}"
			return ast.Bag{
				Pin:  pin,
				Funs: funs,
			}, nil
		}

		f, err := p.bagFun()
		if err != nil {
			return ast.Bag{}, err
		}

		if p.peek.Kind != token.Semicolon {
			return ast.Bag{}, p.unexpected()
		}
		p.advance() // skip ";"

		funs = append(funs, f)
	}
}

func (p *Parser) bagFun() (ast.BagFun, diag.Error) {
	if p.peek.Kind != token.Word {
		return ast.BagFun{}, p.unexpected()
	}
	name := p.word()

	s, err := p.signature()
	if err != nil {
		return ast.BagFun{}, err
	}

	return ast.BagFun{
		Name:      name,
		Signature: s,
	}, nil
}
