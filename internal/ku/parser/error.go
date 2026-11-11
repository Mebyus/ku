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
	// caller should continue parsing from next statement or similar
	// adjacent node
	ssNode ss = iota + 1

	// caller should continue parsing from next top-level node
	ssTop

	// caller should stop parsing immediately due to parsing stop
	ssAbort
)

func (p *Parser) topError(pin sx.Pin, msg string) {
	k := 1
	p.advance()
	for !isTopSync(&p.peek) {
		k += 1
		p.advance()

		if k > 64 {
			p.stop(ast.SyncFail)
			break
		}
	}
	p.report(pin, 0, msg)
}

// report error and sync until after end of current expression
func (p *Parser) syncExp(pin sx.Pin, msg string) (*ast.InvExp, ss) {
	k := 1
	p.advance()

syncLoop:
	for {
		if p.status != 0 {
			p.report(pin, 0, msg)
			return &ast.InvExp{Pin: pin}, ssAbort
		}

		switch p.peek.Kind {
		case token.Semicolon, token.Walrus, token.Assign, token.LeftCurly, token.RightCurly,
			token.Return, token.If, token.Else, token.Const, token.Type:
			// these are likely to terminate expression
			break syncLoop
		case token.Fun:
			// assume new top-level function defenition
			p.report(pin, 0, msg)
			return &ast.InvExp{Pin: pin}, ssTop
		}

		k += 1
		p.advance()

		if k > 64 {
			p.stop(ast.SyncFail)
			p.report(pin, 0, msg)
			return &ast.InvExp{Pin: pin}, ssAbort
		}
	}

	p.report(pin, 0, msg)
	return &ast.InvExp{Pin: pin}, 0
}

// report error and sync until after end of current type specifier
func (p *Parser) syncTypeSpec(pin sx.Pin, msg string) (*ast.InvType, ss) {
	k := 1
	p.advance()

	// sync until type specifier ends
syncLoop:
	for {
		if p.status != 0 {
			p.report(pin, 0, msg)
			return &ast.InvType{Pin: pin}, ssAbort
		}

		switch p.peek.Kind {
		case token.Semicolon, token.Walrus, token.Assign, token.LeftParen,
			token.RightParen, token.Return, token.If, token.Else, token.Const,
			token.True, token.False, token.String, token.Type:
			// these are likely to terminate expression
			// assume it's the end of malformed type specifier
			break syncLoop
		}

		k += 1
		p.advance()

		if k > 64 {
			p.stop(ast.SyncFail)
			p.report(pin, 0, msg)
			return &ast.InvType{Pin: pin}, ssAbort
		}
	}

	p.report(pin, 0, msg)
	return &ast.InvType{Pin: pin}, 0
}

// report error and sync until start of next statement
func (p *Parser) syncNextNode(pin sx.Pin, msg string) ss {
	k := 1
	p.advance()

	// sync until next statement start
syncLoop:
	for {
		if p.status != 0 {
			p.report(pin, 0, msg)
			return ssAbort
		}

		switch p.peek.Kind {
		case token.Semicolon:
			// assume it's the end of malformed statement
			p.advance() // skip ";"
			break syncLoop
		case token.LeftCurly:
			// assume it starts a new block
			break syncLoop
		case token.RightCurly:
			// assume it closes encompassing block
			break syncLoop
		case token.Return, token.If, token.Const, token.Type:
			// these always start a new statement
			break syncLoop
		case token.Fun:
			// assume new top-level function defenition
			p.report(pin, 0, msg)
			return ssTop
		}

		k += 1
		p.advance()

		if k > 64 {
			p.stop(ast.SyncFail)
			p.report(pin, 0, msg)
			return ssAbort
		}
	}

	p.report(pin, 0, msg)
	return ssNode
}

func (p *Parser) report(pin sx.Pin, typ sx.ErrorType, short string, notes ...string) {
	p.pool.Report(pin, typ, short, notes...)
	p.nerr += 1

	if p.nerr > 16 {
		p.stop(ast.LimitReached)
	}
}

func isTopSync(tok *token.Token) bool {
	switch tok.Kind {
	case token.Fun, token.Type, token.EOF:
		return true
	default:
		return false
	}
}
