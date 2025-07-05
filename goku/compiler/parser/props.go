package parser

import (
	"strings"

	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/token"
)

func (p *Parser) gatherProps() diag.Error {
	for {
		if p.c.Kind != token.HashSquare {
			return nil
		}

		props, err := p.PropsBlock()
		if err != nil {
			return err
		}
		p.props = append(p.props, props...)
	}
}

func (p *Parser) PropsBlock() ([]ast.Prop, diag.Error) {
	p.advance() // skip "#["

	var props []ast.Prop
	for {
		if p.c.Kind == token.RightSquare {
			p.advance() // skip "]"
			return props, nil
		}

		prop, err := p.prop()
		if err != nil {
			return nil, err
		}
		props = append(props, prop)
	}
}

func (p *Parser) prop() (ast.Prop, diag.Error) {
	pin := p.c.Pin
	name, err := p.propName()
	if err != nil {
		return ast.Prop{}, err
	}

	if p.c.Kind != token.Assign {
		return ast.Prop{}, p.unexpected()
	}
	p.advance() // skip "="

	exp, err := p.Exp()
	if err != nil {
		return ast.Prop{}, err
	}

	if p.c.Kind != token.Semicolon {
		return ast.Prop{}, p.unexpected()
	}
	p.advance() // skip ";"

	return ast.Prop{
		Exp:  exp,
		Name: name,
		Pin:  pin,
	}, nil
}

func (p *Parser) propName() (string, diag.Error) {
	var parts []string
	part, ok := getPropNamePart(&p.c)
	if !ok {
		return "", p.unexpected()
	}
	p.advance() // skip name part
	parts = append(parts, part)

	for {
		if p.c.Kind != token.Period {
			return strings.Join(parts, "."), nil
		}
		p.advance() // skip "."

		part, ok := getPropNamePart(&p.c)
		if !ok {
			return "", p.unexpected()
		}
		p.advance() // skip name part
		parts = append(parts, part)
	}
}

func getPropNamePart(tok *token.Token) (string, bool) {
	if tok.Kind == token.Word {
		return tok.Data, true
	}
	if tok.IsKeyword() {
		return tok.Kind.String(), true
	}
	return "", false
}

func (p *Parser) takeProps() *[]ast.Prop {
	if len(p.props) == 0 {
		return nil
	}
	props := p.props
	p.props = nil
	return &props
}
