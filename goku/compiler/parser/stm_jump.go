package parser

import (
	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/token"
)

func (p *Parser) Jump() (ast.Statement, diag.Error) {
	pin := p.c.Pin
	p.advance() // skip "jump"

	var s ast.Statement
	switch p.c.Kind {
	case token.LabelNext:
		p.advance() // skip "@.next"
		s = ast.JumpNext{Pin: pin}
	case token.LabelOut:
		p.advance() // skip "@.out"
		s = ast.JumpOut{Pin: pin}
	default:
		return nil, p.unexpected()
	}

	if p.c.Kind != token.Semicolon {
		return nil, p.unexpected()
	}
	p.advance() // skip ";"

	return s, nil
}
