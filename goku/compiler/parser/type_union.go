package parser

import (
	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
)

func (p *Parser) Union() (ast.Union, diag.Error) {
	pin := p.c.Pin
	p.advance() // skip "union"

	fields, err := p.fields()
	if err != nil {
		return ast.Union{}, err
	}

	return ast.Union{
		Fields: fields,
		Pin:    pin,
	}, nil
}
