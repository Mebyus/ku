package parser

import (
	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/token"
)

func (p *Parser) Map() (ast.Map, diag.Error) {
	p.advance() // skip "map"

	if p.peek.Kind != token.LeftParen {
		return ast.Map{}, p.unexpected()
	}
	p.advance() // skip "("

	key, err := p.TypeSpec()
	if err != nil {
		return ast.Map{}, err
	}

	if p.peek.Kind != token.Comma {
		return ast.Map{}, p.unexpected()
	}
	p.advance() // skip ","

	value, err := p.TypeSpec()
	if err != nil {
		return ast.Map{}, err
	}

	if p.peek.Kind != token.RightParen {
		return ast.Map{}, p.unexpected()
	}
	p.advance() // skip ")"

	return ast.Map{
		Key:   key,
		Value: value,
	}, nil
}
