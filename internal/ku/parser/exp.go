package parser

import (
	"fmt"

	"github.com/mebyus/ku/internal/ku/ast"
	"github.com/mebyus/ku/internal/ku/token"
)

func (p *Parser) Exp() ast.Exp {
	return p.operand()
}

func (p *Parser) operand() ast.Operand {
	switch p.peek.Kind {
	case token.Integer:
		tok := p.peek
		p.advance() // skip integer
		return &ast.Integer{
			Pin: tok.Pin,
			Val: tok.Val,
			Aux: tok.Flags,
		}
	default:
		pin := p.peek.Pin
		er := ast.Error{
			Pin:   pin,
			Short: fmt.Sprintf("expected expression or operand, found %s token instead", p.peek.Kind),
		}
		p.addError(&er)
		p.advance() // TODO: we should do a sync here
		return &ast.ErrorExp{Error: er}
	}
}
