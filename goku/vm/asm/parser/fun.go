package parser

import (
	"strings"

	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/vm/asm/ast"
	"github.com/mebyus/ku/goku/vm/tokens"
)

func (p *Parser) topFun() diag.Error {
	p.advance() // skip "#fun"

	pin := p.peek.Pin
	name, err := p.funName()
	if err != nil {
		return err
	}

	atoms, err := p.funBody()
	if err != nil {
		return err
	}

	p.text.Functions = append(p.text.Functions, ast.Fun{
		Atoms: atoms,
		Name:  name,
		Pin:   pin,
	})
	return nil
}

func (p *Parser) funName() (string, diag.Error) {
	if p.peek.Kind != tokens.Word {
		return "", p.unexpected()
	}

	parts := []string{p.peek.Data}
	p.advance() // skip first name part

	for {
		if p.peek.Kind != tokens.Period {
			return strings.Join(parts, "."), nil
		}
		p.advance() // skip "."

		if p.peek.Kind != tokens.Word {
			return "", p.unexpected()
		}
		parts = append(parts, p.peek.Data)
		p.advance() // skip name part
	}
}

func (p *Parser) funBody() ([]ast.Atom, diag.Error) {
	if p.peek.Kind != tokens.LeftCurly {
		return nil, p.unexpected()
	}
	p.advance() // skip "{"

	var atoms []ast.Atom
	for {
		if p.peek.Kind == tokens.RightCurly {
			p.advance() // skip "}"
			return atoms, nil
		}

		atom, err := p.atom()
		if err != nil {
			return nil, err
		}
		atoms = append(atoms, atom)
	}
}
