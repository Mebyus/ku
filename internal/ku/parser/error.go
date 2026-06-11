package parser

import (
	"github.com/mebyus/ku/internal/ku/ast"
	"github.com/mebyus/ku/internal/ku/sx"
	"github.com/mebyus/ku/internal/ku/token"
)

// sync signal after parsing error
//
// most parsing methods in this package return this as second value to
// indicate what happend during subtree parsing and how callers should
// proceed with the result
//
// zero value indicates that caller should proceed with normal path
// other values may indicate that some kind of sync action from parser
// is needed up in the call chain
//
// to clarify: zero value of ss does not mean that subtree parsing went
// without errors or return value is a valid node,
// it merely indicates that caller should proceed as normal
type ss int

const (
	// caller should continue parsing from next statement
	ssNode ss = iota + 1

	// caller should continue parsing from next top-level node
	ssTop

	// caller should stop parsing immediately due to parsing stop
	ssAbort
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
