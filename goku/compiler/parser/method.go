package parser

import (
	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/token"
)

func (p *Parser) Method(traits ast.Traits) diag.Error {
	p.advance() // skip "fun"

	receiver, err := p.Receiver()
	if err != nil {
		return err
	}

	if p.c.Kind != token.Word {
		return p.unexpected()
	}
	name := p.word()

	signature, err := p.signature()
	if err != nil {
		return err
	}

	if p.c.Kind != token.LeftCurly {
		return p.unexpected()
	}

	body, err := p.Block()
	if err != nil {
		return err
	}

	p.text.AddMethod(ast.Method{
		Receiver:  receiver,
		Name:      name,
		Signature: signature,
		Body:      body,
		Traits:    traits,
	})
	return nil
}

func (p *Parser) Receiver() (ast.Receiver, diag.Error) {
	if p.c.Kind != token.LeftParen {
		return ast.Receiver{}, p.unexpected()
	}
	p.advance() // skip "("

	var ptr bool

	if p.c.Kind == token.Asterisk {
		p.advance() // skip "*"
		ptr = true
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
		Ptr:  ptr,
	}, nil
}
