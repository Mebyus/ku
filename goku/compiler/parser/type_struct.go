package parser

import (
	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/token"
)

func (p *Parser) Struct() (ast.Struct, diag.Error) {
	pin := p.peek.Pin
	p.advance() // skip "struct"

	fields, err := p.fields()
	if err != nil {
		return ast.Struct{}, err
	}

	return ast.Struct{
		Fields: fields,
		Pin:    pin,
	}, nil
}

func (p *Parser) fields() ([]ast.Field, diag.Error) {
	if p.peek.Kind != token.LeftCurly {
		return nil, p.unexpected()
	}
	p.advance() // skip "{"

	var fields []ast.Field
	for {
		if p.peek.Kind == token.RightCurly {
			p.advance() // skip "}"
			return fields, nil
		}

		field, err := p.field()
		if err != nil {
			return nil, err
		}
		fields = append(fields, field)

		if p.peek.Kind == token.Comma {
			// Commas are optional between struct fields.
			p.advance() // skip ","
		}
	}
}

func (p *Parser) field() (ast.Field, diag.Error) {
	if p.peek.Kind != token.Word {
		return ast.Field{}, p.unexpected()
	}
	name := p.word()

	if p.peek.Kind != token.Colon {
		return ast.Field{}, p.unexpected()
	}
	p.advance() // consume ":"

	typ, err := p.TypeSpec()
	if err != nil {
		return ast.Field{}, err
	}

	return ast.Field{
		Name: name,
		Type: typ,
	}, nil
}
