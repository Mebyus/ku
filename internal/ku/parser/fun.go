package parser

import (
	"fmt"

	"github.com/mebyus/ku/internal/ku/ast"
	"github.com/mebyus/ku/internal/ku/token"
)

func (p *Parser) Fun() *ast.Fun {
	p.advance() // skip "fun"

	if p.peek.Kind != token.Word {
		p.topError(p.peek.Pin, fmt.Sprintf("expected function name, found %s token instead", p.peek.Kind))
		return nil
	}
	name := p.peek.Data
	pin := p.peek.Pin
	p.advance() // skip function name

	if p.peek.Kind != token.LeftParen {
		p.topError(p.peek.Pin, fmt.Sprintf("expected \"(\" before function param list, found %s token instead", p.peek.Kind))
		return nil
	}
	p.advance() // skip "("

	var params []ast.Param
	for {
		if p.peek.Kind == token.RightParen {
			p.advance() // skip ")"
			break
		}

		param, ok := p.funParam()
		if ok {
			params = append(params, param)
		}

		switch p.peek.Kind {
		case token.Comma:
			p.advance() // skip ","
		case token.RightParen:
			// will be skipped at next iteration
		case token.Word:
			p.report(p.peek.Pin, "missing \",\" between params in function signature")
			// continue parsing, not a serious error
		default:
			p.topError(p.peek.Pin, fmt.Sprintf("expected \")\" or next param in function signature, found %s token instead", p.peek.Kind))
			return nil
		}
	}

	var result ast.TypeSpec
	switch p.peek.Kind {
	case token.RightArrow:
		p.advance() // skip "->"

		if p.peek.Kind == token.LeftCurly {
			p.report(p.peek.Pin, "missing function result type after \"->\"")
		} else {
			typ, ok := p.TypeSpec()
			if !ok {
				return nil
			}
			result = typ
		}
	case token.LeftCurly:
		// function returns nothing
		// continue with parsing function body
	default:
		// try to parse as type specifier
		// maybe it is a slight syntax mistake of missing "->"
		p.report(p.peek.Pin, "missing \"->\" before function result type")

		typ, ok := p.TypeSpec()
		if !ok {
			return nil
		}
		result = typ
	}

	block := p.Block()
	if block == nil {
		return nil
	}

	return &ast.Fun{
		Sig: ast.Signature{
			Result: result,
			Params: params,
		},
		Body: *block,
		Name: name,
		Pin:  pin,
	}
}

func (p *Parser) funParam() (ast.Param, bool) {
	if p.peek.Kind != token.Word {
		p.report(p.peek.Pin, fmt.Sprintf("expected param name, found %s token instead", p.peek.Kind))
		return ast.Param{}, false
	}
	name := p.peek.Data
	pin := p.peek.Pin
	p.advance() // skip param name

	if p.peek.Kind != token.Colon {
		p.report(p.peek.Pin, "missing \":\" before param type specifier")
		// continue parsing, not a serious error
	} else {
		p.advance() // consume ":"
	}

	typ, ok := p.TypeSpec()
	if !ok {
		return ast.Param{}, false
	}

	return ast.Param{
		Name: name,
		Pin:  pin,
		Type: typ,
	}, true
}
