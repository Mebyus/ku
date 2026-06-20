package parser

import (
	"fmt"

	"github.com/mebyus/ku/internal/ku/ast"
	"github.com/mebyus/ku/internal/ku/token"
)

func (p *Parser) topType() {
	p.advance() // skip "type"

	if p.peek.Kind != token.Word {
		p.topError(p.peek.Pin, fmt.Sprintf("expected type name, found %s token instead", &p.peek))
		return
	}

	var def ast.Type

	pin := p.peek.Pin
	name := p.peek.Data
	p.advance() // skip type name

	def.Name = name
	def.Pin = pin

	typ, _ := p.TypeSpec()
	def.Spec = typ
	p.text.AddType(def)

	if p.peek.Kind == token.Semicolon {
		// this semicolon is optional
		p.advance() // skip ";"
	}
}

func (p *Parser) TypeSpec() (ast.TypeSpec, ss) {
	switch p.peek.Kind {
	case token.Word:
		pin := p.peek.Pin
		name := p.peek.Data
		p.advance() // skip word
		return &ast.TypeName{
			Name: name,
			Pin:  pin,
		}, 0
	case token.Struct:
		return p.Struct()
	case token.Ampersand:
		pin := p.peek.Pin
		p.advance() // skip "&"
		typ, s := p.TypeSpec()
		return &ast.Ref{Pin: pin, Type: typ}, s
	case token.LeftSquare:
		if p.next.Kind == token.RightSquare {
			return p.typeSpan()
		}
	}

	return p.syncTypeSpec(p.peek.Pin, fmt.Sprintf("expected type specifier, found %s token instead", &p.peek))
}

func (p *Parser) typeSpan() (ast.TypeSpec, ss) {
	pin := p.peek.Pin

	p.advance() // skip "["
	p.advance() // skip "]"

	typ, s := p.TypeSpec()
	if s != 0 {
		return &ast.InvType{Error: ast.Error{Pin: pin}}, s
	}

	return &ast.Span{
		Pin:  pin,
		Type: typ,
	}, 0
}

func (p *Parser) Struct() (ast.TypeSpec, ss) {
	p.advance() // skip "struct"

	if p.peek.Kind != token.LeftCurly {
		p.report(p.peek.Pin, fmt.Sprintf("expected \"{\" before struct fields, found %s token instead", &p.peek))
		// TODO: sync here
		return &ast.InvType{}, ssNode
	}

	pin := p.peek.Pin
	p.advance() // skip "{"

	r := ast.Struct{Pin: pin}
	for {
		if p.peek.Kind == token.RightCurly {
			p.advance() // skip "}"
			return &r, 0
		}

		f, s := p.field()
		if s != 0 {
			return &r, s
		}
		r.Fields = append(r.Fields, f)

		switch p.peek.Kind {
		case token.Comma:
			p.advance() // skip ","
		case token.Word:
			// continue to next field
		case token.RightCurly:
			// will be skipped at next iteration
		default:
			// TODO: sync + error
			panic("stub")
		}
	}
}

func (p *Parser) field() (ast.Field, ss) {
	if p.peek.Kind != token.Word {
		p.report(p.peek.Pin, fmt.Sprintf("expected field name, found %s token instead", &p.peek))
		return ast.Field{}, ssTop
	}
	name := p.peek.Data
	pin := p.peek.Pin
	p.advance() // skip field name

	if p.peek.Kind != token.Colon {
		p.report(p.peek.Pin, "missing \":\" before field type specifier")
		// continue parsing, not a serious error
	} else {
		p.advance() // consume ":"
	}

	typ, s := p.TypeSpec()
	return ast.Field{
		Name: name,
		Pin:  pin,
		Type: typ,
	}, s
}
