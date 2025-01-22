package parser

import (
	"fmt"

	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/enums/aok"
	"github.com/mebyus/ku/goku/compiler/token"
)

// AssignOrInvoke parses a statement which starts with a word.
func (p *Parser) AssignOrInvoke() (ast.Statement, diag.Error) {
	pack, err := p.Pack()
	if err != nil {
		return nil, err
	}

	switch p.c.Kind {
	case token.Semicolon:
		return p.invoke(pack)
	case token.AddAssign, token.SubAssign, token.MulAssign, token.DivAssign, token.RemAssign:
	case token.Assign, token.Walrus:
	default:
		return nil, p.unexpected()
	}

	k, ok := aok.FromToken(p.c.Kind)
	if !ok {
		panic(fmt.Sprintf("\"%s\" token must be assign operator", p.c.Kind))
	}
	op := ast.AssignOp{Pin: p.c.Pin, Kind: k}
	p.advance() // skip assign operator

	value, err := p.Pack()
	if err != nil {
		return nil, err
	}

	if p.c.Kind != token.Semicolon {
		return nil, p.unexpected()
	}
	p.advance() // skip ";"

	return ast.Assign{
		Op:     op,
		Target: pack,
		Value:  value,
	}, nil
}

func (p *Parser) invoke(exp ast.Exp) (ast.Invoke, diag.Error) {
	call, ok := exp.(ast.Call)
	if ok {
		p.advance() // skip ";"
		return ast.Invoke{Call: call}, nil
	}

	return ast.Invoke{}, &diag.SimpleMessageError{
		Pin:  exp.Span().Pin,
		Text: fmt.Sprintf("%s expression used as a statement, only call expression can be used this way", exp.Kind()),
	}
}
