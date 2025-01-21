package parser

import (
	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/token"
)

// Pack parses a list of non-pack expressions. List must contain at least one
// such expression.
//
// If list has only one element then return that element as expression.
// Returns pack expression from that list otherwise.
func (p *Parser) Pack() (ast.Exp, diag.Error) {
	var list []ast.Exp
	for {
		if isPackEnd(p.c.Kind) {
			switch len(list) {
			case 0:
				return nil, p.unexpected()
			case 1:
				return list[0], nil
			default:
				return ast.Pack{List: list}, nil
			}
		}

		exp, err := p.Exp()
		if err != nil {
			return nil, err
		}
		list = append(list, exp)

		if p.c.Kind == token.Comma {
			p.advance() // skip ","
		} else if isPackEnd(p.c.Kind) {
			// will be skipped at next iteration
		} else {
			return nil, p.unexpected()
		}
	}
}

// Pack expression should be terminated by one of these: ";", "=", ":=".
func isPackEnd(k token.Kind) bool {
	return k == token.Semicolon || k == token.Assign || k == token.Walrus
}
