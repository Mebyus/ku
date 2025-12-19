package parser

import (
	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/enums/aok"
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
		if isPackEnd(p.peek.Kind) {
			switch len(list) {
			case 0:
				return nil, p.unexpected()
			case 1:
				return list[0], nil
			default:
				return ast.Pack{List: list}, nil
			}
		}

		exp, err := p.packExp()
		if err != nil {
			return nil, err
		}
		list = append(list, exp)

		if p.peek.Kind == token.Comma {
			p.advance() // skip ","
		} else if isPackEnd(p.peek.Kind) {
			// will be skipped at next iteration
		} else {
			return nil, p.unexpected()
		}
	}
}

// parses expression that can appear as pack element
func (p *Parser) packExp() (ast.Exp, diag.Error) {
	if p.peek.Kind == token.Underscore {
		pin := p.peek.Pin
		p.advance()
		return ast.Blank{Pin: pin}, nil
	}

	return p.Exp()
}

// Pack expression should be terminated by one of these: ";", "=", ":=".
func isPackEnd(k token.Kind) bool {
	if k == token.Semicolon {
		return true
	}
	_, ok := aok.FromToken(k)
	return ok
}
