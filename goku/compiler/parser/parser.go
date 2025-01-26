package parser

import (
	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/lexer"
	"github.com/mebyus/ku/goku/compiler/token"
)

type Parser struct {
	text ast.Text

	lx lexer.Stream

	c token.Token
	n token.Token

	props []ast.Prop
}

func (p *Parser) Nodes() (*ast.Text, diag.Error) {
	err := p.parse()
	if err != nil {
		return nil, err
	}
	return &p.text, nil
}

func (p *Parser) Text() (*ast.Text, diag.Error) {
	_, err := p.Build()
	if err != nil {
		return nil, err
	}
	_, err = p.ImportBlocks()
	if err != nil {
		return nil, err
	}
	return p.Nodes()
}

func FromStream(stream lexer.Stream) *Parser {
	p := Parser{lx: stream}
	p.init()
	return &p
}

func ParseStream(stream lexer.Stream) (*ast.Text, diag.Error) {
	return FromStream(stream).Text()
}
