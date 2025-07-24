package parser

import (
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/kub/ast"
	"github.com/mebyus/ku/goku/kub/token"
)

func (p *Parser) pkgParse() diag.Error {
	for {
		if p.peek.Kind == token.EOF {
			return nil
		}

		err := p.pkgTop()
		if err != nil {
			return err
		}
	}
}

func (p *Parser) pkgTop() diag.Error {
	switch p.peek.Kind {
	case token.Set:
		return p.setBlock()
	case token.Module:
		return p.module()
	default:
		return p.unexpected()
	}
}

func (p *Parser) setBlock() diag.Error {
	p.advance() // skip "set"

	if p.peek.Kind != token.LeftCurly {
		return p.unexpected()
	}
	p.advance() // skip "{"

	for {
		switch p.peek.Kind {
		case token.RightCurly:
			p.advance() // skip "}"
			return nil
		default:
			set, err := p.set()
			if err != nil {
				return err
			}
			p.pkg.Sets = append(p.pkg.Sets, set)
		}
	}
}

func (p *Parser) namePart() (ast.Word, diag.Error) {
	var val string
	switch p.peek.Kind {
	case token.Word:
		val = p.peek.Data
	case token.Main, token.Link:
		val = p.peek.Kind.String()
	default:
		return ast.Word{}, p.unexpected()
	}
	pin := p.peek.Pin
	p.advance() // skip name part

	return ast.Word{
		Str: val,
		Pin: pin,
	}, nil
}

func (p *Parser) set() (ast.Set, diag.Error) {
	var parts []ast.Word

	part, err := p.namePart()
	if err != nil {
		return ast.Set{}, err
	}
	parts = append(parts, part)

	for {
		if p.peek.Kind != token.Period {
			break
		}
		p.advance() // skip "."

		part, err := p.namePart()
		if err != nil {
			return ast.Set{}, err
		}
		parts = append(parts, part)
	}

	if p.peek.Kind != token.Assign {
		return ast.Set{}, p.unexpected()
	}
	p.advance() // skip "="

	if p.peek.Kind != token.String {
		return ast.Set{}, p.unexpected()
	}
	exp := ast.String{
		Val: p.peek.Data,
		Pin: p.peek.Pin,
	}
	p.advance() // skip exp string

	if p.peek.Kind != token.Semicolon {
		return ast.Set{}, p.unexpected()
	}
	p.advance() // skip ";"

	return ast.Set{
		Name: ast.Name{Parts: parts},
		Exp:  exp,
	}, nil
}

func (p *Parser) module() diag.Error {
	p.advance() // skip "module"

	if p.peek.Kind != token.String {
		return p.unexpected()
	}
	name := ast.String{
		Val: p.peek.Data,
		Pin: p.peek.Pin,
	}
	if name.Val == "" {
		return &diag.SimpleMessageError{
			Pin:  name.Pin,
			Text: "empty module name",
		}
	}
	p.advance() // skip module name

	if p.peek.Kind != token.LeftCurly {
		return p.unexpected()
	}
	p.advance() // skip "{"

	module := ast.Module{Name: name}
	for {
		switch p.peek.Kind {
		case token.Main:
			main, err := p.mainEntry()
			if err != nil {
				return err
			}
			if module.Main != nil {
				return &diag.SimpleMessageError{
					Pin:  main.Pin,
					Text: "duplicate main unit inside module",
				}
			}
			module.Main = &main
		case token.Unit:
			unit, err := p.unitEntry()
			if err != nil {
				return err
			}
			module.Units = append(module.Units, unit)
		case token.Link:
			link, err := p.linkEntry()
			if err != nil {
				return err
			}
			module.Links = append(module.Links, link)
		case token.RightCurly:
			p.advance() // skip "}"
			p.pkg.Modules = append(p.pkg.Modules, module)
			return nil
		default:
			return p.unexpected()
		}
	}
}

func (p *Parser) unitEntry() (ast.UnitEntry, diag.Error) {
	p.advance() // skip "unit"

	if p.peek.Kind != token.String {
		return ast.UnitEntry{}, p.unexpected()
	}
	pin := p.peek.Pin
	val := p.peek.Data
	if val == "" {
		return ast.UnitEntry{}, &diag.SimpleMessageError{
			Pin:  pin,
			Text: "empty unit path",
		}
	}
	p.advance() // skip string

	if p.peek.Kind != token.Semicolon {
		return ast.UnitEntry{}, p.unexpected()
	}
	p.advance() // skip ";"

	return ast.UnitEntry{
		Val: val,
		Pin: pin,
	}, nil
}

func (p *Parser) mainEntry() (ast.MainEntry, diag.Error) {
	p.advance() // skip "main"

	if p.peek.Kind != token.String {
		return ast.MainEntry{}, p.unexpected()
	}
	pin := p.peek.Pin
	val := p.peek.Data
	if val == "" {
		return ast.MainEntry{}, &diag.SimpleMessageError{
			Pin:  pin,
			Text: "empty main path",
		}
	}
	p.advance() // skip string

	if p.peek.Kind != token.Semicolon {
		return ast.MainEntry{}, p.unexpected()
	}
	p.advance() // skip ";"

	return ast.MainEntry{
		Val: val,
		Pin: pin,
	}, nil
}

func (p *Parser) linkEntry() (ast.LinkEntry, diag.Error) {
	p.advance() // skip "link"

	if p.peek.Kind != token.String {
		return ast.LinkEntry{}, p.unexpected()
	}
	pin := p.peek.Pin
	val := p.peek.Data
	p.advance() // skip string

	if p.peek.Kind != token.Semicolon {
		return ast.LinkEntry{}, p.unexpected()
	}
	p.advance() // skip ";"

	return ast.LinkEntry{
		Val: val,
		Pin: pin,
	}, nil
}
