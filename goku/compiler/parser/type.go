package parser

import (
	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/token"
)

func (p *Parser) topType(traits ast.Traits) diag.Error {
	t, err := p.Type(traits)
	if err != nil {
		return err
	}
	p.text.AddType(t)
	return nil
}

func (p *Parser) Type(traits ast.Traits) (ast.Type, diag.Error) {
	p.advance() // skip "type"

	if p.c.Kind != token.Word {
		return ast.Type{}, p.unexpected()
	}
	name := p.word()

	var bags []ast.Word
	if p.c.Kind == token.In {
		p.advance() // skip "in"
		list, err := p.bagList()
		if err != nil {
			return ast.Type{}, err
		}
		bags = list
	}

	if p.c.Kind != token.RightArrow {
		return ast.Type{}, p.unexpected()
	}
	p.advance() // skip "=>"

	spec, err := p.CustomTypeSpec()
	if err != nil {
		return ast.Type{}, err
	}

	return ast.Type{
		Name:   name,
		Bags:   bags,
		Spec:   spec,
		Traits: traits,
	}, nil
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
	case token.Type:
		return p.AnyType(), nil
	case token.Fun:
		return p.FunType()
	default:
		return nil, p.unexpected()
	}
}

// ResultTypeSpec parses type specifier in function signature return type.
//
// This variant includes tuples and forms.
func (p *Parser) ResultTypeSpec() (ast.TypeSpec, diag.Error) {
	if p.c.Kind == token.LeftParen {
		return p.TupleOrForm()
	}
	return p.TypeSpec()
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
	case token.Union:
		return p.Union()
	case token.Struct:
		return p.Struct()
	case token.LeftCurly:
		if p.n.Kind == token.RightCurly {
			pin := p.c.Pin
			p.advance() // skip "{"
			p.advance() // skip "}"
			return ast.Trivial{Pin: pin}, nil
		}
	case token.Bag:
		return p.Bag()
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

func (p *Parser) AnyType() ast.AnyType {
	pin := p.c.Pin

	p.advance() // skip "type"

	return ast.AnyType{Pin: pin}
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

func (p *Parser) bagList() ([]ast.Word, diag.Error) {
	if p.c.Kind != token.LeftParen {
		return nil, p.unexpected()
	}
	p.advance() // skip "("

	var list []ast.Word
	for {
		if p.c.Kind == token.RightParen {
			p.advance() // skip ")"
			return list, nil
		}

		if p.c.Kind != token.Word {
			return nil, p.unexpected()
		}
		name := p.word()
		list = append(list, name)

		if p.c.Kind == token.Comma {
			p.advance() // skip ","
		} else if p.c.Kind == token.RightParen {
			// will be skipped at next iteration
		} else {
			return nil, p.unexpected()
		}
	}
}
