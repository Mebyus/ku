package parser

import (
	"fmt"

	"github.com/mebyus/ku/internal/ku/ast"
	"github.com/mebyus/ku/internal/ku/token"
)

func (p *Parser) Chain() (ast.Operand, ss) {
	pin := p.peek.Pin
	name := p.peek.Data
	if p.next.Kind == token.AlterZero {
		p.advance() // skip word (symbol name)
		p.advance() // skip ".{}"
		return &ast.SymZeroExp{
			Pin:  pin,
			Name: name,
		}, 0
	}

	chain := ast.Chain{
		Pin:  pin,
		Name: name,
	}
	p.advance() // skip word (symbol name)

	for {
		switch p.peek.Kind {
		case token.Address:
			p.advance() // skip ".&"
			if len(chain.Parts) == 0 {
				return &ast.SymAddr{
					Name: chain.Name,
					Pin:  chain.Pin,
				}, 0
			}
			panic("stub")
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
				p.report(p.peek.Pin, fmt.Sprintf("expected name after select operand, found %s token instead", &p.peek))
				p.advance() // TODO: sync
			}
		case token.LeftParen:
			args, s := p.args()
			if s != 0 {
				// return &ast.InvExp{Error: }, 0
			}
			if len(chain.Parts) == 0 {
				return &ast.SymCall{
					Args: args,
					Name: chain.Name,
					Pin:  chain.Pin,
				}, 0
			}
			panic("stub")
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

func (p *Parser) args() ([]ast.Exp, ss) {
	p.advance() // skip "("

	var args []ast.Exp
	for {
		if p.peek.Kind == token.RightParen {
			p.advance() // skip ")"
			return args, 0
		}

		exp, s := p.Exp()
		if s != 0 {
			return args, s
		}
		args = append(args, exp)

		switch p.peek.Kind {
		case token.Comma:
			p.advance() // skip ","
		case token.RightParen:
			// will be skipped at next iteration
		default:
			p.report(p.peek.Pin, fmt.Sprintf("unexpected %s token after argument expression", &p.peek))
			return args, ssNode
		}
	}
}
