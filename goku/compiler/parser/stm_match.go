package parser

import (
	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/token"
)

func (p *Parser) Match(exp ast.Exp) (ast.Match, diag.Error) {
	m := ast.Match{Exp: exp}
	for {
		switch p.c.Kind {
		case token.RightArrow:
			c, err := p.matchCase()
			if err != nil {
				return ast.Match{}, err
			}
			m.Cases = append(m.Cases, c)
		case token.Else:
			p.advance() // skip "else"

			body, err := p.Block()
			if err != nil {
				return ast.Match{}, err
			}
			m.Else = &body
			return m, nil
		default:
			return m, nil
		}
	}
}

func (p *Parser) matchCase() (ast.MatchCase, diag.Error) {
	p.advance() // skip "=>"

	list, err := p.ExpList()
	if err != nil {
		return ast.MatchCase{}, err
	}
	if len(list) == 0 {
		return ast.MatchCase{}, &diag.SimpleMessageError{
			Pin:  p.c.Pin,
			Text: "case with no expressions",
		}
	}

	block, err := p.Block()
	if err != nil {
		return ast.MatchCase{}, err
	}

	return ast.MatchCase{
		List: list,
		Body: block,
	}, nil
}

func (p *Parser) ExpList() ([]ast.Exp, diag.Error) {
	var list []ast.Exp
	for {
		if p.c.Kind == token.LeftCurly {
			return list, nil
		}

		exp, err := p.Exp()
		if err != nil {
			return nil, err
		}
		list = append(list, exp)

		if p.c.Kind == token.Comma {
			p.advance() // skip ","
		} else if p.c.Kind == token.LeftCurly {
			// will cause return at next iteration
		} else {
			return nil, p.unexpected()
		}
	}
}
