package parser

import (
	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/token"
)

func (p *Parser) RegBag() diag.Error {
	p.advance() // skip "bag"

	if p.peek.Kind != token.Word {
		return p.unexpected()
	}
	name := p.word()

	if p.peek.Kind != token.Colon {
		return p.unexpected()
	}
	p.advance() // skip ":"

	typ, err := p.TypeSpec()
	if err != nil {
		return err
	}

	if p.peek.Kind != token.RightArrow {
		return p.unexpected()
	}
	p.advance() // skip "=>"

	if p.peek.Kind != token.Word {
		return p.unexpected()
	}
	bag := p.word()

	tab, err := p.bagTab()
	if err != nil {
		return err
	}

	reg := ast.RegBag{
		Type:    typ,
		Name:    name,
		BagName: bag,
		Tab:     tab,
	}
	p.text.AddRegBag(reg)
	return nil
}

func (p *Parser) bagTab() ([]ast.BagTabField, diag.Error) {
	if p.peek.Kind != token.LeftCurly {
		return nil, p.unexpected()
	}
	p.advance() // skip "{"

	var tab []ast.BagTabField
	for {
		if p.peek.Kind == token.RightCurly {
			p.advance() // skip "}"
			return tab, nil
		}

		if p.peek.Kind != token.Word {
			return nil, p.unexpected()
		}
		name := p.word()

		if p.peek.Kind != token.Colon {
			return nil, p.unexpected()
		}
		p.advance() // skip ":"

		if p.peek.Kind != token.Word {
			return nil, p.unexpected()
		}
		fun := p.word()

		if p.peek.Kind != token.Comma {
			return nil, p.unexpected()
		}
		p.advance() // skip ","

		tab = append(tab, ast.BagTabField{
			Name: name,
			Fun:  fun,
		})
	}
}
