package parser

import (
	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/token"
)

func (p *Parser) Enum() (ast.Enum, diag.Error) {
	base := p.TypeName()

	p.advance() // skip "{"

	var entries []ast.EnumEntry
	for {
		if p.peek.Kind == token.RightCurly {
			p.advance() // skip "}"

			return ast.Enum{
				Base:    base,
				Entries: entries,
			}, nil
		}

		entry, err := p.EnumEntry()
		if err != nil {
			return ast.Enum{}, err
		}
		entries = append(entries, entry)

		if p.peek.Kind == token.Comma {
			// Commas are optional between enum entries.
			p.advance() // skip ","
		}
	}
}

func (p *Parser) EnumEntry() (ast.EnumEntry, diag.Error) {
	if p.peek.Kind != token.Word {
		return ast.EnumEntry{}, p.unexpected()
	}
	name := p.word()

	if p.peek.Kind != token.Assign {
		// entry without explicitly assigned value
		return ast.EnumEntry{Name: name}, nil
	}

	p.advance() // skip "="

	exp, err := p.Exp()
	if err != nil {
		return ast.EnumEntry{}, err
	}

	return ast.EnumEntry{
		Name: name,
		Exp:  exp,
	}, nil
}
