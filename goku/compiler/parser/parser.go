package parser

import (
	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/cerp"
	"github.com/mebyus/ku/goku/compiler/lexer"
	"github.com/mebyus/ku/goku/compiler/token"
)

type Parser struct {
	text ast.Text

	s lexer.Stream

	c token.Token
	n token.Token
}

func (p *Parser) Imports() ([]ast.ImportBlock, cerp.Error) {
	return nil, nil
}

func (p *Parser) Nodes() (*ast.Text, cerp.Error) {
	return &p.text, nil
}

func (p *Parser) Text() (*ast.Text, cerp.Error) {
	_, err := p.Imports()
	if err != nil {
		return nil, err
	}
	return p.Nodes()
}

func FromStream(stream lexer.Stream) *Parser {
	p := Parser{s: stream}
	p.init()
	return &p
}

func ParseStream(stream lexer.Stream) (*ast.Text, cerp.Error) {
	return FromStream(stream).Text()
}
