package parser

import (
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/sm"
	"github.com/mebyus/ku/goku/vm/asm/ast"
	"github.com/mebyus/ku/goku/vm/lexer"
	"github.com/mebyus/ku/goku/vm/tokens"
)

type Parser struct {
	text ast.Text

	lx *lexer.Lexer

	peek tokens.Token
	next tokens.Token
}

func (p *Parser) Text() (*ast.Text, diag.Error) {
	err := p.parse()
	if err != nil {
		return nil, err
	}
	return &p.text, nil
}

func FromText(text *sm.Text) *Parser {
	p := Parser{lx: lexer.FromText(text)}
	p.init()
	return &p
}
