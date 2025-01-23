package parser

import (
	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/token"
)

// InitExp parses expression in extended form, which is allowed for init expressions.
func (p *Parser) InitExp() (ast.Exp, diag.Error) {
	if p.c.Kind == token.Quest {
		pin := p.c.Pin
		p.advance()
		return ast.Dirty{Pin: pin}, nil
	}

	return p.Exp()
}
