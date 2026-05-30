package parser

import (
	"github.com/mebyus/ku/internal/ku/ast"
	"github.com/mebyus/ku/internal/ku/sx"
	"github.com/mebyus/ku/internal/ku/token"
)

func (p *Parser) topError(pin sx.Pin, msg string) {
	er := ast.Error{
		Short: msg,
		Pin:   pin,
	}

	er.Tokens = append(er.Tokens, p.peek)
	p.advance()

	for !isTopSync(&p.peek) {
		er.Tokens = append(er.Tokens, p.peek)
		p.advance()

		if len(er.Tokens) > 64 {
			p.abort(ast.ErrorSyncFailed)
			break
		}
	}
	p.addError(&er)
}

func (p *Parser) addError(er *ast.Error) {
	p.text.Errors = append(p.text.Errors, er)
	if len(p.text.Errors) > 16 {
		p.abort(ast.ErrorLimitReached)
	}
}

func (p *Parser) report(pin sx.Pin, msg string) {
	p.addError(&ast.Error{
		Short: msg,
		Pin:   pin,
	})
}

func isTopSync(tok *token.Token) bool {
	switch tok.Kind {
	case token.Fun, token.EOF:
		return true
	default:
		return false
	}
}
