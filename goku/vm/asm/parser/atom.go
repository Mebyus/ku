package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/srcmap"
	"github.com/mebyus/ku/goku/vm/asm/ast"
	"github.com/mebyus/ku/goku/vm/opc"
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
	case tokens.Reg:
		return p.register()
	case tokens.Word:
		return p.symbol()
	case tokens.DecInteger:
		return p.decInteger()
	case tokens.Label:
		return p.label()
	default:
		return nil, p.unexpected()
	}
}

func (p *Parser) register() (ast.Register, diag.Error) {
	name := p.peek.Data
	pin := p.peek.Pin
	p.advance() // skip register name

	reg, err := parseRegisterName(pin, name)
	if err != nil {
		return ast.Register{}, err
	}

	return ast.Register{
		Name: reg,
		Pin:  pin,
	}, nil
}

func (p *Parser) label() (ast.Label, diag.Error) {
	name := p.peek.Data
	pin := p.peek.Pin
	p.advance() // skip label name

	return ast.Label{
		Name: name,
		Pin:  pin,
	}, nil
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

func unknownRegister(pin srcmap.Pin, s string) diag.Error {
	return &diag.SimpleMessageError{
		Pin:  pin,
		Text: fmt.Sprintf("unknown register \"%s\"", s),
	}
}

func parseRegisterName(pin srcmap.Pin, s string) (opc.Register, diag.Error) {
	switch s {
	case "sp":
		return opc.RegFP, nil
	case "ip":
		return opc.RegIP, nil
	case "sc":
		return opc.RegSC, nil
	case "fp":
		return opc.RegFP, nil
	case "clock":
		return opc.RegClock, nil
	case "cf":
		return opc.RegCF, nil
	default:
		if !strings.HasPrefix(s, "r") {
			return 0, unknownRegister(pin, s)
		}
		n, err := strconv.ParseUint(s[1:], 10, 64)
		if err != nil {
			return 0, unknownRegister(pin, s)
		}
		if n >= 64 {
			return 0, unknownRegister(pin, s)
		}
		return opc.Register(n), nil
	}
}
