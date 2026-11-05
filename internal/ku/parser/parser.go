package parser

import (
	"github.com/mebyus/ku/internal/ku/ast"
	"github.com/mebyus/ku/internal/ku/lexer"
	"github.com/mebyus/ku/internal/ku/sx"
	"github.com/mebyus/ku/internal/ku/token"
)

type Parser struct {
	text ast.Text

	lx *lexer.Lexer

	// used for reporting errors
	pool *sx.Pool

	// token at current parser position
	peek token.Token

	// next token after peek
	next token.Token

	// number of skipped valid tokens (not invalid ones) due
	// to syncs after errors
	nskip uint32

	// number of reported errors
	nerr uint32

	// number of invalid tokens encountered
	ninv uint32

	// parsing should be stopped when this is not 0
	status ast.Status
}

func FromText(pool *sx.Pool, text *sx.Text) *Parser {
	p := Parser{
		pool: pool,
		lx:   lexer.FromText(text),
	}
	p.init()
	return &p
}

func ParseText(pool *sx.Pool, text *sx.Text) *ast.Text {
	return FromText(pool, text).Parse()
}

func (p *Parser) Parse() *ast.Text {
	p.parse()
	return &p.text
}

func (p *Parser) advance() {
	if p.status != 0 {
		setEOF(&p.peek, uint64(p.status))
		return
	}

	var tok token.Token
	p.lx.Lex(&tok)
	if tok.Kind == token.INV {
		// TODO: consume invalid tokens with loop?
		p.ninv += 1
		if p.ninv > 32 {
			p.stop(ast.InvalidTokens)
			setEOF(&p.peek, uint64(p.status))
			return
		}
	}

	p.peek = p.next
	p.next = tok
}

// stop parsing with specified status/error
func (p *Parser) stop(status ast.Status) {
	if p.status != 0 {
		// ignore subsequent status changes after first
		return
	}

	p.status = status
}

func (p *Parser) init() {
	p.advance()
	p.advance()
}

func setEOF(tok *token.Token, val uint64) {
	tok.Val = val
	tok.Kind = token.EOF
	tok.Data = ""
	tok.Flags = 0
}
