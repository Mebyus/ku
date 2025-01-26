package parser

import (
	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/token"
)

// Build parses build block from source text. Returns nil if there is no build block.
func (p *Parser) Build() (*ast.Build, diag.Error) {
	if p.c.Kind != token.Build {
		return nil, nil
	}

	p.advance() // skip "#build"

	body, err := p.Block()
	if err != nil {
		return nil, err
	}
	build := &ast.Build{Body: body}
	p.text.Build = build
	return build, nil
}
