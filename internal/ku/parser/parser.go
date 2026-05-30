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

	// token at current parser position
	peek token.Token

	// next token after peek
	next token.Token

	// total number of illegal tokens encountered
	illegals uint32

	// parsing should stop when this is set to true
	stop bool
}

func FromText(text *sx.Text) *Parser {
	p := Parser{lx: lexer.FromText(text)}
	p.init()
	return &p
}

func ParseText(text *sx.Text) *ast.Text {
	return FromText(text).Parse()
}

func (p *Parser) Parse() *ast.Text {
	p.parse()
	return &p.text
}

func (p *Parser) advance() {
	if p.stop {
		setEOF(&p.peek, 1) // TODO: specify EOF error code
		return
	}

	var tok token.Token
	p.lx.Lex(&tok)
	if tok.Kind == token.ILG {
		p.illegals += 1
		if p.illegals > 32 {
			setEOF(&p.peek, 1) // TODO: specify EOF error code
			p.abort(ast.ErrorIllegalTokens)
			return
		}
	}

	p.peek = p.next
	p.next = tok
}

// stop parsing with specified status/error
func (p *Parser) abort(status ast.Status) {
	if p.stop {
		// ignore subsequent aborts after first
		return
	}

	p.text.Status = status
	p.stop = true
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
