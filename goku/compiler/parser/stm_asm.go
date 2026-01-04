package parser

import (
	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/token"
)

func (p *Parser) Asm() (*ast.Asm, diag.Error) {
	p.advance() // skip "asm"

	var a ast.Asm
	if p.peek.Kind != token.Word {
		return nil, p.unexpected()
	}
	a.Arch = p.word()

	err := p.asmArgs(&a)
	if err != nil {
		return nil, err
	}

	if p.peek.Kind != token.LeftCurly {
		return nil, p.unexpected()
	}
	p.advance() // skip "{"
	for {
		if p.peek.Kind == token.RightCurly {
			p.advance() // skip "}"
			return &a, nil
		}

		switch p.peek.Kind {
		case token.Word:
			m := p.word() // asm instruction mnemonic

			var rest []token.Token
			for p.peek.Kind != token.Semicolon {
				rest = append(rest, p.peek)
				p.advance()
			}
			p.advance() // skip ";"

			a.Nodes = append(a.Nodes, ast.AsmInst{
				Rest:     rest,
				Mnemonic: m.Str,
				Pin:      m.Pin,
			})
		default:
			return nil, p.unexpected()
		}
	}
}

func (p *Parser) asmArgs(a *ast.Asm) diag.Error {
	if p.peek.Kind != token.LeftParen {
		return p.unexpected()
	}
	p.advance() // skip "("

	for {
		if p.peek.Kind == token.RightParen {
			p.advance() // skip ")"
			return nil
		}

		switch {
		case p.peek.Kind == token.Word && p.next.Kind == token.Assign:
			name := p.word() // register name
			p.advance()      // skip "="
			if p.peek.Kind != token.Word {
				return p.unexpected()
			}
			symbol := p.word() // init symbol name
			a.Inits = append(a.Inits, ast.RegInit{Name: name, Symbol: symbol})
		case p.peek.Kind == token.Word && p.next.Kind == token.RightArrow:
			name := p.word() // register name
			p.advance()      // skip "->"
			if p.peek.Kind != token.Word {
				return p.unexpected()
			}
			symbol := p.word() // output symbol name
			a.Outs = append(a.Outs, ast.RegOut{Name: name, Symbol: symbol})
		case p.peek.Kind == token.Quest && p.next.Kind == token.RightArrow:
			p.advance() // skip "?"
			p.advance() // skip "->"
			if p.peek.Kind != token.Word {
				return p.unexpected()
			}
			name := p.word() // clobber name
			a.Clobbers = append(a.Clobbers, ast.Clobber{Name: name})
		default:
			return p.unexpected()
		}

		if p.peek.Kind != token.Semicolon {
			return p.unexpected()
		}
		p.advance() // skip ";"
	}
}
