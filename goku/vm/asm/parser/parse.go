package parser

import (
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/sm"
	"github.com/mebyus/ku/goku/vm/asm/ast"
	"github.com/mebyus/ku/goku/vm/tokens"
)

func Parse(text *sm.Text) (*ast.Text, diag.Error) {
	p := FromText(text)
	return p.Text()
}

func (p *Parser) parse() diag.Error {
	for {
		if p.peek.Kind == tokens.EOF {
			return nil
		}

		err := p.top()
		if err != nil {
			return err
		}
	}
}

func (p *Parser) top() diag.Error {
	switch p.peek.Kind {
	case tokens.Fun:
		return p.topFun()
	case tokens.Data:
		panic("not implemented")
	case tokens.Entry:
		return p.topEntry()
	default:
		return p.unexpected()
	}
}

func (p *Parser) topEntry() diag.Error {
	if p.text.Entry.Name != "" {
		return &diag.SimpleMessageError{
			Pin:  p.peek.Pin,
			Text: "entrypoint was already declared in program",
		}
	}

	p.advance() // skip "#entry"

	if p.peek.Kind != tokens.Word {
		return p.unexpected()
	}
	name := p.peek.Data
	pin := p.peek.Pin
	p.advance() // skip entrypoint name

	if p.peek.Kind != tokens.Semicolon {
		return p.unexpected()
	}
	p.advance() // skip ";"

	p.text.Entry = ast.Entry{
		Name: name,
		Pin:  pin,
	}
	return nil
}
