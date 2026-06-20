package parser

import (
	"fmt"

	"github.com/mebyus/ku/internal/ku/ast"
	"github.com/mebyus/ku/internal/ku/token"
)

func (p *Parser) topFun() {
	p.advance() // skip "fun"

	if p.peek.Kind != token.Word {
		p.topError(p.peek.Pin, fmt.Sprintf("expected function name, found %s token instead", &p.peek))
		return
	}

	var fun ast.Fun

	name := p.peek.Data
	pin := p.peek.Pin
	p.advance() // skip function name

	fun.Name = name
	fun.Pin = pin

	if p.peek.Kind != token.LeftParen {
		p.topError(p.peek.Pin, fmt.Sprintf("expected \"(\" before function param list, found %s token instead", &p.peek))
		p.text.AddFun(fun)
		return
	}
	p.advance() // skip "("

	for {
		if p.peek.Kind == token.RightParen {
			p.advance() // skip ")"
			break
		}

		param, s := p.funParam()
		switch s {
		case ssTop, ssAbort:
			p.text.AddFun(fun)
			return
		}
		fun.Sig.Params = append(fun.Sig.Params, param)

		switch p.peek.Kind {
		case token.Comma:
			p.advance() // skip ","
		case token.RightParen:
			// will be skipped at next iteration
		case token.Word:
			p.report(p.peek.Pin, "missing \",\" between params in function signature")
			// continue parsing, not a serious error
		default:
			p.topError(p.peek.Pin, fmt.Sprintf("expected \")\" or next param in function signature, found %s token instead", &p.peek))
			p.text.AddFun(fun)
			return
		}
	}

	switch p.peek.Kind {
	case token.RightArrow:
		p.advance() // skip "->"

		if p.peek.Kind == token.LeftCurly {
			p.report(p.peek.Pin, "missing function result type after \"->\"")
		} else {
			typ, s := p.TypeSpec()
			if s != 0 {
				p.text.AddFun(fun)
				return
			}
			fun.Sig.Result = typ
		}
	case token.LeftCurly:
		// function returns nothing
		// continue with parsing function body
	case token.Semicolon:
		// function forward declaration
		p.advance() // skip ";"
		p.text.AddStub(ast.FunStub{
			Sig:  fun.Sig,
			Name: fun.Name,
			Pin:  fun.Pin,
		})
		return
	default:
		// try to parse as type specifier
		// maybe it is a slight syntax mistake of missing "->"
		p.report(p.peek.Pin, "missing \"->\" before function result type")

		typ, s := p.TypeSpec()
		if s != 0 {
			p.text.AddFun(fun)
			return
		}
		fun.Sig.Result = typ
	}

	switch p.peek.Kind {
	case token.LeftCurly:
		p.block(&fun.Body)
		p.skipBadCurly()
		p.text.AddFun(fun)
		return
	case token.Semicolon:
		// function forward declaration
		p.advance() // skip ";"
		p.text.AddStub(ast.FunStub{
			Sig:  fun.Sig,
			Name: fun.Name,
			Pin:  fun.Pin,
		})
		return
	default:
		p.report(p.peek.Pin, fmt.Sprintf("expected \"{\" as function body start or \";\" as the end of function declaration, found %s token instead", &p.peek))

		// assume it's forward declaration with missing ";"
		p.text.AddStub(ast.FunStub{
			Sig:  fun.Sig,
			Name: fun.Name,
			Pin:  fun.Pin,
		})
		// syncronize
	}
}

func (p *Parser) funParam() (ast.Param, ss) {
	if p.peek.Kind != token.Word {
		p.report(p.peek.Pin, fmt.Sprintf("expected param name, found %s token instead", &p.peek))
		return ast.Param{}, ssTop
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

	typ, s := p.TypeSpec()
	return ast.Param{
		Name: name,
		Pin:  pin,
		Type: typ,
	}, s
}

func (p *Parser) skipBadCurly() {
	switch p.peek.Kind {
	case token.LeftCurly, token.RightCurly:
		// continue execution
	default:
		return
	}

	er := ast.Error{
		Pin:   p.peek.Pin,
		Short: "unbalanced curly brace",
	}
	er.Tokens = append(er.Tokens, p.peek)
	p.advance()

	for {
		switch p.peek.Kind {
		case token.LeftCurly, token.RightCurly:
			er.Tokens = append(er.Tokens, p.peek)
			p.advance()
		default:
			p.addError(&er)
			return
		}
	}
}
