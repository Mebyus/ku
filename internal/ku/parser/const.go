package parser

import (
	"fmt"

	"github.com/mebyus/ku/internal/ku/ast"
	"github.com/mebyus/ku/internal/ku/token"
)

// accepts
//
//	const name = 1;
//	const name := 1;
//	const name : = 1;
//	const name: uint = 1;
func (p *Parser) topConst() {
	p.advance() // skip "fun"

	if p.peek.Kind != token.Word {
		p.topError(p.peek.Pin, fmt.Sprintf("expected constant name, found %s token instead", &p.peek))
		return
	}

	var c ast.Const

	name := p.peek.Data
	pin := p.peek.Pin
	p.advance() // skip constant name

	c.Name = name
	c.Pin = pin

	if p.peek.Kind == token.Colon {
		p.advance() // skip ":"

		if p.peek.Kind == token.Assign {
			p.advance() // skip "="
		} else {
			typ, s := p.TypeSpec()
			if s != 0 {
				p.text.AddConst(c)
				return
			}
			c.Type = typ

			if p.peek.Kind != token.Assign {
				p.topError(p.peek.Pin, fmt.Sprintf("expected \"=\" before constant init expression, found %s token instead", &p.peek))
				p.text.AddConst(c)
				return
			}
			p.advance() // skip "="
		}
	} else if p.peek.Kind == token.Walrus {
		p.advance() // skip ":="
	} else if p.peek.Kind == token.Assign {
		p.advance() // skip "="
	} else {
		p.topError(p.peek.Pin, fmt.Sprintf("expected \"=\" before constant init expression, found %s token instead", &p.peek))
		p.text.AddConst(c)
		return
	}

	exp, _ := p.Exp()
	c.Exp = exp

	if p.peek.Kind != token.Semicolon {
		p.report(p.peek.Pin, "missing \";\" after constant init expression")
	} else {
		p.advance() // skip ";"
	}

	p.text.AddConst(c)
}
