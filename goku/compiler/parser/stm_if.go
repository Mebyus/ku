package parser

import (
	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/token"
)

func (p *Parser) If() (ast.Statement, diag.Error) {
	// switch p.n.Kind {
	// case token.RightArrow, token.Else:
	// 	return p.matchBool()
	// }

	p.advance() // skip "if"

	exp, err := p.Exp()
	if err != nil {
		return nil, err
	}

	switch p.c.Kind {
	case token.RightArrow, token.Else:
		return p.Match(exp)
	case token.LeftCurly:
		// continue regular if statement
	default:
		return nil, p.unexpected()
	}

	body, err := p.Block()
	if err != nil {
		return nil, err
	}

	var elseIfs []ast.IfClause
	for {
		if p.c.Kind == token.Else && p.n.Kind == token.If {
			p.advance() // skip "else"
		} else {
			break
		}

		clause, err := p.ifClause()
		if err != nil {
			return nil, err
		}
		elseIfs = append(elseIfs, clause)
	}

	var elseBody *ast.Block
	if p.c.Kind == token.Else {
		p.advance() // skip "else"

		var body ast.Block
		body, err = p.Block()
		if err != nil {
			return nil, err
		}
		elseBody = &body
	}

	return ast.If{
		If: ast.IfClause{
			Exp:  exp,
			Body: body,
		},
		ElseIfs: elseIfs,
		Else:    elseBody,
	}, nil
}

func (p *Parser) ifClause() (ast.IfClause, diag.Error) {
	p.advance() // skip "if"

	exp, err := p.Exp()
	if err != nil {
		return ast.IfClause{}, err
	}
	if p.c.Kind != token.LeftCurly {
		return ast.IfClause{}, p.unexpected()
	}
	body, err := p.Block()
	if err != nil {
		return ast.IfClause{}, err
	}

	return ast.IfClause{
		Exp:  exp,
		Body: body,
	}, nil
}
