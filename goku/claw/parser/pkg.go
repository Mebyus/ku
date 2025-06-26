package parser

import (
	"github.com/mebyus/ku/goku/claw/ast"
	"github.com/mebyus/ku/goku/claw/token"
	"github.com/mebyus/ku/goku/compiler/diag"
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
		return p.set()
	case token.Module:
		return p.module()
	default:
		return p.unexpected()
	}
}

func (p *Parser) set() diag.Error {
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
			p.advance() // TODO: remove dummy code
		}
	}
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
