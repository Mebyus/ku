package parser

import (
	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/srcmap"
	"github.com/mebyus/ku/goku/compiler/token"
)

func (p *Parser) TupleOrForm() (ast.TypeSpec, diag.Error) {
	pin := p.peek.Pin
	p.advance() // skip "("

	if p.peek.Kind == token.Word && p.next.Kind == token.Colon {
		return p.form()
	}

	return p.tuple(pin)
}

func (p *Parser) tuple(pin srcmap.Pin) (ast.Tuple, diag.Error) {
	var types []ast.TypeSpec
	for {
		if p.peek.Kind == token.RightParen {
			p.advance() // skip ")"
			return ast.Tuple{
				Pin:   pin,
				Types: types,
			}, nil
		}

		typ, err := p.TypeSpec()
		if err != nil {
			return ast.Tuple{}, err
		}
		types = append(types, typ)

		if p.peek.Kind == token.Comma {
			p.advance() // skip ","
		} else if p.peek.Kind == token.RightParen {
			// will be skipped at next iteration
		} else {
			return ast.Tuple{}, p.unexpected()
		}
	}
}

func (p *Parser) form() (ast.Form, diag.Error) {
	var fields []ast.Field
	for {
		if p.peek.Kind == token.RightParen {
			p.advance() // skip ")"
			return ast.Form{Fields: fields}, nil
		}

		field, err := p.field()
		if err != nil {
			return ast.Form{}, err
		}
		fields = append(fields, field)

		if p.peek.Kind == token.Comma {
			p.advance() // skip ","
		} else if p.peek.Kind == token.RightParen {
			// will be skipped at next iteration
		} else {
			return ast.Form{}, p.unexpected()
		}
	}
}
