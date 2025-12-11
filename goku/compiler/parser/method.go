package parser

import (
	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/token"
)

func (p *Parser) topMethod(traits ast.Traits) diag.Error {
	m, err := p.Method(traits)
	if err != nil {
		return err
	}
	p.text.AddMethod(m)
	return nil
}

func (p *Parser) Method(traits ast.Traits) (ast.Method, diag.Error) {
	p.advance() // skip "fun"

	receiver, err := p.Receiver()
	if err != nil {
		return ast.Method{}, err
	}

	err = p.unsafe(&traits)
	if err != nil {
		return ast.Method{}, err
	}

	f, err := p.fun()
	if err != nil {
		return ast.Method{}, err
	}

	return ast.Method{
		Receiver:  receiver,
		Name:      f.Name,
		Signature: f.Signature,
		Body:      f.Body,
		Traits:    traits,
	}, nil
}

func (p *Parser) Receiver() (ast.Receiver, diag.Error) {
	if p.c.Kind != token.LeftParen {
		return ast.Receiver{}, p.unexpected()
	}
	p.advance() // skip "("

	var kind ast.ReceiverKind
	switch p.c.Kind {
	case token.Asterisk:
		p.advance() // skip "*"
		kind = ast.ReceiverPtr
	case token.Ampersand:
		p.advance() // skip "&"
		kind = ast.ReceiverRef
	default:
		kind = ast.ReceiverVal
	}

	if p.c.Kind != token.Word {
		return ast.Receiver{}, p.unexpected()
	}
	name := p.word()

	if p.c.Kind != token.RightParen {
		return ast.Receiver{}, p.unexpected()
	}
	p.advance() // skip ")"

	return ast.Receiver{
		Name: name,
		Kind: kind,
	}, nil
}
