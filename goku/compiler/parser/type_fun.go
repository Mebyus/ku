package parser

import (
	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
)

func (p *Parser) FunType() (ast.FunType, diag.Error) {
	p.advance() // skip "fun"
	s, err := p.signature()
	if err != nil {
		return ast.FunType{}, err
	}
	return ast.FunType{Signature: s}, nil
}
