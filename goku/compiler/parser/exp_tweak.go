package parser

import (
	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/token"
)

func (p *Parser) Object() (ast.Object, diag.Error) {
	pin := p.c.Pin
	p.advance() // skip "{"

	var fields []ast.ObjField
	for {
		if p.c.Kind == token.RightCurly {
			p.advance() // skip "}"
			return ast.Object{
				Pin:    pin,
				Fields: fields,
			}, nil
		}

		field, err := p.objField()
		if err != nil {
			return ast.Object{}, err
		}
		fields = append(fields, field)

		if p.c.Kind == token.Comma {
			p.advance() // skip ","
		} else if p.c.Kind == token.RightCurly {
			// will be skipped at next iteration
		} else {
			return ast.Object{}, p.unexpected()
		}
	}
}

func (p *Parser) tweak(chain ast.Chain) (ast.Tweak, diag.Error) {
	p.advance() // skip ".{"

	var fields []ast.ObjField
	for {
		if p.c.Kind == token.RightCurly {
			p.advance() // skip "}"
			return ast.Tweak{
				Chain:  chain,
				Fields: fields,
			}, nil
		}

		field, err := p.objField()
		if err != nil {
			return ast.Tweak{}, err
		}
		fields = append(fields, field)

		if p.c.Kind == token.Comma {
			p.advance() // skip ","
		} else if p.c.Kind == token.RightCurly {
			// will be skipped at next iteration
		} else {
			return ast.Tweak{}, p.unexpected()
		}
	}
}

func (p *Parser) objField() (ast.ObjField, diag.Error) {
	if p.c.Kind != token.Word {
		return ast.ObjField{}, p.unexpected()
	}
	name := p.word()

	if p.c.Kind != token.Colon {
		return ast.ObjField{}, p.unexpected()
	}
	p.advance() // consume ":"

	exp, err := p.Exp()
	if err != nil {
		return ast.ObjField{}, err
	}

	return ast.ObjField{
		Name: name,
		Exp:  exp,
	}, nil
}
