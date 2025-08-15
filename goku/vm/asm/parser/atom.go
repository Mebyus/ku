package parser

import (
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/vm/asm/ast"
	"github.com/mebyus/ku/goku/vm/tokens"
)

func (p *Parser) atom() (ast.Atom, diag.Error) {
	switch p.peek.Kind {
	case tokens.Word:
		return p.instruction()
	case tokens.Label:
		return p.place()
	default:
		return nil, p.unexpected()
	}
}

func (p *Parser) place() (ast.Place, diag.Error) {
	name := p.peek.Data
	pin := p.peek.Pin
	p.advance() // skip label name

	if p.peek.Kind != tokens.Colon {
		return ast.Place{}, p.unexpected()
	}
	p.advance() // skip ":"

	return ast.Place{
		Name: name,
		Pin:  pin,
	}, nil
}

func (p *Parser) instruction() (ast.Instruction, diag.Error) {
	mnemonic := p.peek.Data
	pin := p.peek.Pin
	p.advance() // skip instruction mnemonic

	var operands []ast.Operand
	for {
		if p.peek.Kind == tokens.Semicolon {
			p.advance() // skip ";"
			return ast.Instruction{
				Operands: operands,
				Mnemonic: mnemonic,
				Pin:      pin,
			}, nil
		}

		operand, err := p.operand()
		if err != nil {
			return ast.Instruction{}, err
		}

		if p.peek.Kind == tokens.Comma {
			p.advance() // skip ","
		} else if p.peek.Kind == tokens.Semicolon {
			// will be skipped at next iteration
		} else {
			return ast.Instruction{}, p.unexpected()
		}

		operands = append(operands, operand)
	}
}

func (p *Parser) operand() (ast.Operand, diag.Error) {
	switch p.peek.Kind {
	case tokens.Word:
		return p.symbol()
	case tokens.DecInteger:
		return p.decInteger()
	default:
		return nil, p.unexpected()
	}
}

func (p *Parser) symbol() (ast.Symbol, diag.Error) {
	name := p.peek.Data
	pin := p.peek.Pin
	p.advance() // skip symbol name

	return ast.Symbol{
		Name: name,
		Pin:  pin,
	}, nil
}

func (p *Parser) decInteger() (ast.Integer, diag.Error) {
	val := p.peek.Val
	pin := p.peek.Pin
	p.advance() // skip integer

	return ast.Integer{
		Val: val,
		Pin: pin,
	}, nil
}
