package parser

import (
	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/token"
)

func (p *Parser) Tuple() (ast.Tuple, diag.Error) {
	pin := p.c.Pin
	p.advance() // skip "("

	var types []ast.TypeSpec
	for {
		if p.c.Kind == token.RightParen {
			p.advance() // skip ")"
			return ast.Tuple{
				Pin:   pin,
				Types: types,
			}, nil
		}

		typ, err := p.TypeSpec()
		if err != nil {
			return ast.Tuple{}, err
		}
		types = append(types, typ)

		if p.c.Kind == token.Comma {
			p.advance() // skip ","
		} else if p.c.Kind == token.RightParen {
			// will be skipped at next iteration
		} else {
			return ast.Tuple{}, p.unexpected()
		}
	}
}
