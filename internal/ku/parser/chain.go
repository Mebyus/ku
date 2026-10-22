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
		case token.SliceAll:
			pin := p.peek.Pin
			p.advance() // skip "[:]"
			return &ast.SliceAll{
				Pin:   pin,
				Chain: chain,
			}, 0
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
				panic("stub")
			}
			if len(chain.Parts) == 0 {
				return &ast.SymCall{
					Args: args,
					Name: chain.Name,
					Pin:  chain.Pin,
				}, 0
			}
			panic("stub")
		case token.LeftSquare:
			pin := p.peek.Pin
			x, s := p.index()
			if s != 0 {
				panic("stub")
			}
			if x.slice {
				// TODO: should return finished expression here
				panic("stub")
			}
			chain.Parts = append(chain.Parts, &ast.Index{
				Pin: pin,
				Exp: x.a,
			})
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

type indexlike struct {
	a ast.Exp

	// not nil only when slice part was parsed
	b ast.Exp

	// true when slice part was parsed instead of index
	slice bool
}

func (p *Parser) index() (indexlike, ss) {
	p.advance() // skip "["
	if p.peek.Kind == token.Colon {
		// TODO: parse slice
		panic("stub")
	}

	exp, s := p.Exp()
	if s != 0 {
		return indexlike{}, s
	}

	if p.peek.Kind != token.RightSquare {
		return indexlike{}, ssNode
	}
	p.advance() // skip "]"

	return indexlike{a: exp}, 0
}
