package parser

import (
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/sm"
	"github.com/mebyus/ku/goku/kub/ast"
	"github.com/mebyus/ku/goku/kub/lexer"
	"github.com/mebyus/ku/goku/kub/token"
)

type Parser struct {
	pkg  ast.Package
	unit ast.Unit

	lx lexer.Stream

	peek token.Token
	next token.Token
}

func (p *Parser) Package() (*ast.Package, diag.Error) {
	err := p.pkgParse()
	if err != nil {
		return nil, err
	}
	return &p.pkg, nil
}

func (p *Parser) Unit() (*ast.Unit, diag.Error) {
	err := p.unitParse()
	if err != nil {
		return nil, err
	}
	return &p.unit, nil
}

func FromStream(stream lexer.Stream) *Parser {
	p := Parser{lx: stream}
	p.init()
	return &p
}

func FromText(text *sm.Text) *Parser {
	return FromStream(lexer.FromText(text))
}
