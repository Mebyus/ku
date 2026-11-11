package parser

import (
	"fmt"
	"strings"

	"github.com/mebyus/ku/internal/ku/ast"
	"github.com/mebyus/ku/internal/ku/sx"
	"github.com/mebyus/ku/internal/ku/token"
)

// parse charms for unit-level symbols and store them
// for later
func (p *Parser) topCharms() {
	p.resetCharms()
	for {
		switch p.peek.Kind {
		case token.Pub:
			if p.pub {
				p.report(p.peek.Pin, 0, "duplicate usage of \"pub\" keyword")
			} else {
				p.pub = true
			}
			p.advance() // skip "pub"
		case token.HashSquare:
			s := p.charm()
			if s != 0 {
				return
			}
		default:
			return
		}
	}
}

func (p *Parser) charm() ss {
	p.advance() // skip "#["

	var name string
	var pin sx.Pin
	switch p.peek.Kind {
	case token.RightSquare:
		p.report(p.peek.Pin, 0, "empty charm")
		p.advance() // skip "]"
		return 0
	case token.Word:
		pin = p.peek.Pin
		name = p.peek.Data
		p.advance() // skip first word in name
	default:
		return ssNode
	}

	if p.peek.Kind == token.Period {
		// preallocate 2 slots since we already have first part
		// of the name and after "." should come at least one more
		parts := make([]string, 0, 2)
		parts = append(parts, name)

	scanParts:
		for {
			if p.peek.Kind != token.Period {
				break scanParts
			} else {
				p.advance() // skip "."
			}

			switch p.peek.Kind {
			case token.Word:
				parts = append(parts, p.peek.Data)
				p.advance() // skip word
			case token.LeftParen, token.RightSquare:
				p.report(p.peek.Pin, 0, "missing charm name part")
				break scanParts
			case token.Period:
				p.report(p.peek.Pin, 0, "missing charm name part")
			default:
				p.report(p.peek.Pin, 0, fmt.Sprintf("expected charm name part, found %s token instead", &p.peek))
				return ssNode
			}
		}

		name = strings.Join(parts, ".")
	}

	var args []ast.Arg
	if p.peek.Kind == token.LeftParen {
		p.advance() // skip "("

		for {
			if p.peek.Kind == token.RightParen {
				if len(args) == 0 {
					p.report(p.peek.Pin, 0, "no charm arguments")
				}
				p.advance() // skip ")"
				break
			}

			switch p.peek.Kind {
			case token.Word:
				args = append(args, ast.Arg{
					Text: p.peek.Data,
					Pin:  p.peek.Pin,
					Type: ast.ArgWord,
				})
				p.advance() // skip word
			case token.String:
				args = append(args, ast.Arg{
					Text: p.peek.Data,
					Pin:  p.peek.Pin,
					Type: ast.ArgStr,
				})
				p.advance() // skip string
			case token.Integer:
				args = append(args, ast.Arg{
					Int:  p.peek.Val,
					Pin:  p.peek.Pin,
					Type: ast.ArgInt,
				})
				p.advance() // skip integer
			case token.Comma:
				p.report(p.peek.Pin, 0, "empty charm argument")
			default:
				p.report(p.peek.Pin, 0, fmt.Sprintf("expected charm argument, found %s token instead", &p.peek))
				return ssNode
			}

			switch p.peek.Kind {
			case token.RightParen:
				// continue execution, will be consumed on next iteration
			case token.Comma:
				p.advance() // skip ","
			case token.RightSquare:
				p.report(p.peek.Pin, 0, "missing \")\" after charm arguments")
				p.advance() // skip "]"
				p.charms = append(p.charms, ast.Charm{
					Args: args,
					Name: name,
					Pin:  pin,
				})
				return 0
			case token.Word, token.String, token.Integer:
				p.report(p.peek.Pin, 0, "missing \",\" before charm argument")
				// proceed to next argument
			default:
				p.report(p.peek.Pin, 0, fmt.Sprintf("expected \")\" or \",\" before next charm argument, found %s token instead", &p.peek))
				return ssNode
			}
		}
	}

	if p.peek.Kind != token.RightSquare {
		p.report(p.peek.Pin, 0, "missing \"]\" to close a charm")
	} else {
		p.advance() // skip "]"
	}

	p.charms = append(p.charms, ast.Charm{
		Args: args,
		Name: name,
		Pin:  pin,
	})
	return 0
}

func (p *Parser) resetCharms() {
	p.charms = nil
	p.pub = false
}
