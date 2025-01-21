package parser

import (
	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/token"
)

func (p *Parser) Type(traits ast.Traits) diag.Error {
	p.advance() // skip "type"

	if p.c.Kind != token.Word {
		return p.unexpected()
	}
	name := p.word()

	if p.c.Kind != token.RightArrow {
		return p.unexpected()
	}
	p.advance() // skip "=>"

	spec, err := p.CustomTypeSpec()
	if err != nil {
		return err
	}

	p.text.AddType(ast.Type{
		Name:   name,
		Spec:   spec,
		Traits: traits,
	})
	return nil
}

// TypeSpec parses type specifier in regular usage form.
//
// It is a restricted set of type specifiers allowed in:
//
//   - function param
//   - field definition
//   - variable and constant definition
//   - inside another type specifier
func (p *Parser) TypeSpec() (ast.TypeSpec, diag.Error) {
	switch p.c.Kind {
	case token.Word:
		if p.n.Kind == token.Period {
			return p.TypeFullName()
		}
		return p.TypeName(), nil
	case token.Asterisk:
		if p.n.Kind == token.Any {
			return p.AnyPointer(), nil
		}
		return p.Pointer()
	case token.ArrayPointer:
		return p.ArrayPointer()
	case token.Chunk:
		return p.Chunk()
	case token.LeftSquare:
		return p.Array()
	default:
		return nil, p.unexpected()
	}
}

// CustomTypeSpec parses type specifier in extended form.
//
// It includes all forms allowed in regular form as well those
// only allowed in custom type definition.
func (p *Parser) CustomTypeSpec() (ast.TypeSpec, diag.Error) {
	switch p.c.Kind {
	case token.Word:
		if p.n.Kind == token.LeftCurly {
			return p.Enum()
		}
	case token.Struct:
	}
	return p.TypeSpec()
}

func (p *Parser) Array() (ast.Array, diag.Error) {
	p.advance() // skip "["

	size, err := p.Exp()
	if err != nil {
		return ast.Array{}, err
	}

	if p.c.Kind != token.RightSquare {
		return ast.Array{}, p.unexpected()
	}
	p.advance() // skip "]"

	t, err := p.TypeSpec()
	if err != nil {
		return ast.Array{}, err
	}

	return ast.Array{
		Size: size,
		Type: t,
	}, nil
}

func (p *Parser) ArrayPointer() (ast.ArrayPointer, diag.Error) {
	p.advance() // skip "[*]"

	t, err := p.TypeSpec()
	if err != nil {
		return ast.ArrayPointer{}, err
	}

	return ast.ArrayPointer{Type: t}, nil
}

func (p *Parser) Chunk() (ast.Chunk, diag.Error) {
	p.advance() // skip "[]"

	t, err := p.TypeSpec()
	if err != nil {
		return ast.Chunk{}, err
	}

	return ast.Chunk{Type: t}, nil
}

func (p *Parser) Pointer() (ast.Pointer, diag.Error) {
	p.advance() // skip "*"

	t, err := p.TypeSpec()
	if err != nil {
		return ast.Pointer{}, err
	}

	return ast.Pointer{Type: t}, nil
}

func (p *Parser) AnyPointer() ast.AnyPointer {
	pin := p.c.Pin

	p.advance() // skip "*"
	p.advance() // skip "any"

	return ast.AnyPointer{Pin: pin}
}

func (p *Parser) TypeFullName() (ast.TypeFullName, diag.Error) {
	iname := p.word()
	p.advance() // skip "."

	if p.c.Kind != token.Word {
		return ast.TypeFullName{}, p.unexpected()
	}
	name := p.word()

	return ast.TypeFullName{
		Import: iname,
		Name:   name,
	}, nil
}

func (p *Parser) TypeName() ast.TypeName {
	name := p.word()
	return ast.TypeName{Name: name}
}
