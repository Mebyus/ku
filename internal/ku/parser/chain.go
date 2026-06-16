package parser

import (
	"fmt"

	"github.com/mebyus/ku/internal/ku/ast"
	"github.com/mebyus/ku/internal/ku/token"
)

func (p *Parser) Chain() (ast.Operand, ss) {
	chain := ast.Chain{
		Pin:  p.peek.Pin,
		Name: p.peek.Data,
	}
	p.advance() // skip word (symbol name)

	for {
		switch p.peek.Kind {
		case token.Period:
			p.advance() // skip "."

			switch p.peek.Kind {
			case token.Word:
				pin := p.peek.Pin
				name := p.peek.Data
				p.advance() // skip select name
				chain.Parts = append(chain.Parts, &ast.Select{
					Name: name,
					Pin:  pin,
				})
			default:
				p.report(p.peek.Pin, fmt.Sprintf("expected name after select operand, found %s token instead", p.peek.Kind))
				p.advance() // TODO: sync
			}
		default:
			if len(chain.Parts) == 0 {
				return &ast.SymExp{
					Pin:  chain.Pin,
					Name: chain.Name,
				}, 0
			}
			return &chain, 0
		}
	}
}
