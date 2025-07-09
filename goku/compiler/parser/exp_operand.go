package parser

import (
	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/token"
)

func (p *Parser) Operand() (ast.Operand, diag.Error) {
	switch p.c.Kind {
	case token.BinInteger:
		tok := p.c
		p.advance()
		return ast.Integer{
			Pin: tok.Pin,
			Val: tok.Val,
			Aux: uint32(ast.IntBin),
		}, nil
	case token.OctInteger:
		tok := p.c
		p.advance()
		return ast.Integer{
			Pin: tok.Pin,
			Val: tok.Val,
			Aux: uint32(ast.IntOct),
		}, nil
	case token.DecInteger:
		tok := p.c
		p.advance()
		return ast.Integer{
			Pin: tok.Pin,
			Val: tok.Val,
			Aux: uint32(ast.IntDec),
		}, nil
	case token.HexInteger:
		tok := p.c
		p.advance()
		return ast.Integer{
			Pin: tok.Pin,
			Val: tok.Val,
			Aux: uint32(ast.IntHex),
		}, nil
	case token.DecFloat:
		tok := p.c
		p.advance()
		return ast.Float{
			Pin: tok.Pin,
			Val: tok.Data,
		}, nil
	case token.String:
		tok := p.c
		p.advance() // skip string
		return ast.String{
			Val: tok.Data,
			Pin: tok.Pin,
		}, nil
	case token.Rune:
		tok := p.c
		p.advance() // skip rune
		return ast.Rune{
			Val: tok.Val,
			Pin: tok.Pin,
		}, nil
	case token.True:
		pin := p.c.Pin
		p.advance() // skip "true"
		return ast.True{Pin: pin}, nil
	case token.False:
		pin := p.c.Pin
		p.advance() // skip "false"
		return ast.False{Pin: pin}, nil
	case token.Nil:
		pin := p.c.Pin
		p.advance() // skip "nil"
		return ast.Nil{Pin: pin}, nil
	case token.TypeId:
		return p.TypeId()
	case token.ErrorId:
		return p.ErrorId()
	case token.Size:
		return p.Size()
	case token.Cast:
		return p.Cast()
	case token.Tint:
		return p.Tint()
	// case token.MemCast:
	// 	return p.memcast()
	case token.LeftCurly:
		return p.Object()
	case token.Word, token.Unsafe:
		return p.Chain()
	case token.LeftParen:
		return p.Paren()
	case token.LeftSquare, token.Chunk:
		return p.List()
	// case token.Chunk:
	// 	return p.chunkStartOperand()
	case token.Period:
		return p.DotName()
	default:
		return nil, p.unexpected()
	}
}

func (p *Parser) Paren() (ast.Paren, diag.Error) {
	p.advance() // skip "("
	exp, err := p.Exp()
	if err != nil {
		return ast.Paren{}, err
	}
	if p.c.Kind != token.RightParen {
		return ast.Paren{}, err
	}
	p.advance() // skip ")"
	return ast.Paren{Exp: exp}, nil
}

func (p *Parser) Chain() (ast.Operand, diag.Error) {
	var chain ast.Chain
	switch p.c.Kind {
	case token.Word:
		start := p.word()
		chain = ast.Chain{Start: start}
	case token.Unsafe:
		unsafe, err := p.Unsafe()
		if err != nil {
			return nil, err
		}
		chain = ast.Chain{Parts: []ast.Part{unsafe}}
	case token.Test:
		test, err := p.SelectTest()
		if err != nil {
			return nil, err
		}
		chain = ast.Chain{Parts: []ast.Part{test}}
	default:
		return nil, p.unexpected()
	}

	for {
		var err diag.Error
		var part ast.Part

		switch p.c.Kind {
		case token.LeftParen:
			return p.call(chain)
		case token.Tweak:
			return p.tweak(chain)
		case token.Period:
			p.advance() // skip "."

			switch p.c.Kind {
			case token.Test:
				part, err = p.SelectTest()
			case token.Unsafe:
				part, err = p.Unsafe()
			case token.Word:
				part = p.Select()
			default:
				return nil, p.unexpected()
			}
		case token.DerefSelect:
			part, err = p.DerefSelect()
		case token.Deref:
			part = p.Deref()
		case token.Address:
			return p.ref(chain), nil
		case token.DerefIndex:
			part, err = p.DerefIndex()
		// case token.BagSelect:
		// 	part, err = p.bagSelectPart()
		case token.LeftSquare:
			var s SliceOrIndex
			s, err = p.SliceOrIndex()
			if err != nil {
				return nil, err
			}
			if !s.Index {
				return ast.Slice{
					Chain: chain,
					Start: s.Exp,
					End:   s.End,
				}, nil
			}
			part = ast.Index{Exp: s.Exp}
		default:
			if len(chain.Parts) == 0 {
				return ast.Symbol{
					Name: chain.Start.Str,
					Pin:  chain.Start.Pin,
				}, nil
			}
			return chain, nil
		}
		if err != nil {
			return nil, err
		}
		chain.Parts = append(chain.Parts, part)
	}
}

func (p *Parser) ref(chain ast.Chain) ast.Ref {
	p.advance() // skip ".&"
	return ast.Ref{Chain: chain}
}

func (p *Parser) DotName() (ast.DotName, diag.Error) {
	pin := p.c.Pin
	p.advance() // skip "."

	if p.c.Kind != token.Word {
		return ast.DotName{}, p.unexpected()
	}
	name := p.word()

	return ast.DotName{
		Pin:  pin,
		Name: name.Str,
	}, nil
}

func (p *Parser) Deref() ast.Deref {
	pin := p.c.Pin
	p.advance() // skip ".*"
	return ast.Deref{Pin: pin}
}

func (p *Parser) DerefSelect() (ast.DerefSelect, diag.Error) {
	p.advance() // skip ".*."

	if p.c.Kind != token.Word {
		return ast.DerefSelect{}, p.unexpected()
	}
	name := p.word()

	return ast.DerefSelect{Name: name}, nil
}

func (p *Parser) DerefIndex() (ast.DerefIndex, diag.Error) {
	p.advance() // skip ".["

	exp, err := p.Exp()
	if err != nil {
		return ast.DerefIndex{}, err
	}
	if p.c.Kind != token.RightSquare {
		return ast.DerefIndex{}, p.unexpected()
	}
	p.advance() // skip "]"

	return ast.DerefIndex{Exp: exp}, nil
}

func (p *Parser) Select() ast.Select {
	name := p.word()
	return ast.Select{Name: name}
}

func (p *Parser) SelectTest() (ast.SelectTest, diag.Error) {
	p.advance() // skip "test"

	if p.c.Kind != token.Period {
		return ast.SelectTest{}, p.unexpected()
	}
	p.advance() // skip "."

	if p.c.Kind != token.Word {
		return ast.SelectTest{}, p.unexpected()
	}

	name := p.word()
	return ast.SelectTest{Name: name}, nil
}

func (p *Parser) Unsafe() (ast.Unsafe, diag.Error) {
	pin := p.c.Pin

	p.advance() // skip "unsafe"

	if p.c.Kind != token.Period {
		return ast.Unsafe{}, p.unexpected()
	}
	p.advance() // skip "."

	if p.c.Kind != token.Word {
		return ast.Unsafe{}, p.unexpected()
	}

	name := p.word()
	return ast.Unsafe{
		Pin:  pin,
		Name: name.Str,
	}, nil
}

func (p *Parser) call(chain ast.Chain) (ast.Call, diag.Error) {
	args, err := p.Args()
	if err != nil {
		return ast.Call{}, err
	}

	return ast.Call{
		Chain: chain,
		Args:  args,
	}, nil
}

func (p *Parser) Args() ([]ast.Exp, diag.Error) {
	p.advance() // skip "("

	var args []ast.Exp
	for {
		if p.c.Kind == token.RightParen {
			p.advance() // skip ")"
			return args, nil
		}

		exp, err := p.Exp()
		if err != nil {
			return nil, err
		}
		args = append(args, exp)

		if p.c.Kind == token.Comma {
			p.advance() // skip ","
		} else if p.c.Kind == token.RightParen {
			// will be skipped at next iteration
		} else {
			return nil, p.unexpected()
		}
	}
}

type SliceOrIndex struct {
	// Index expression (when Index = true) or start expression.
	Exp ast.Exp

	// Valid only when field Index = false.
	End ast.Exp

	// True when this struct carries index expression.
	Index bool
}

func (p *Parser) SliceOrIndex() (SliceOrIndex, diag.Error) {
	p.advance() // skip "["

	if p.c.Kind == token.Colon {
		p.advance() // skip ":"
		if p.c.Kind == token.RightSquare {
			p.advance() // skip "]"
			return SliceOrIndex{}, nil
		}

		end, err := p.Exp()
		if err != nil {
			return SliceOrIndex{}, err
		}
		if p.c.Kind != token.RightSquare {
			return SliceOrIndex{}, p.unexpected()
		}
		p.advance() // skip "]"
		return SliceOrIndex{End: end}, nil
	}

	exp, err := p.Exp()
	if err != nil {
		return SliceOrIndex{}, err
	}
	if p.c.Kind == token.Colon {
		p.advance() // skip ":"
		if p.c.Kind == token.RightSquare {
			p.advance() // skip "]"
			return SliceOrIndex{Exp: exp}, nil
		}
		end, err := p.Exp()
		if err != nil {
			return SliceOrIndex{}, err
		}
		if p.c.Kind != token.RightSquare {
			return SliceOrIndex{}, p.unexpected()
		}
		p.advance() // skip "]"
		return SliceOrIndex{
			Exp: exp,
			End: end,
		}, nil
	}

	if p.c.Kind != token.RightSquare {
		return SliceOrIndex{}, p.unexpected()
	}
	p.advance() // skip "]"
	return SliceOrIndex{
		Exp:   exp,
		Index: true,
	}, nil
}

func (p *Parser) TypeId() (ast.TypeId, diag.Error) {
	p.advance() // skip "#typeid"

	if p.c.Kind != token.LeftParen {
		return ast.TypeId{}, p.unexpected()
	}
	p.advance() // skip "("

	if p.c.Kind != token.Word {
		return ast.TypeId{}, p.unexpected()
	}
	name := p.word()

	if p.c.Kind != token.RightParen {
		return ast.TypeId{}, p.unexpected()
	}
	p.advance() // skip ")"

	return ast.TypeId{Name: name}, nil
}

func (p *Parser) ErrorId() (ast.ErrorId, diag.Error) {
	p.advance() // skip "#error"

	if p.c.Kind != token.LeftParen {
		return ast.ErrorId{}, p.unexpected()
	}
	p.advance() // skip "("

	if p.c.Kind != token.Word {
		return ast.ErrorId{}, p.unexpected()
	}
	name := p.word()

	if p.c.Kind != token.RightParen {
		return ast.ErrorId{}, p.unexpected()
	}
	p.advance() // skip ")"

	return ast.ErrorId{Name: name}, nil
}

func (p *Parser) Cast() (ast.Cast, diag.Error) {
	p.advance() // skip "#size"

	if p.c.Kind != token.LeftParen {
		return ast.Cast{}, p.unexpected()
	}
	p.advance() // skip "("

	typ, err := p.TypeSpec()
	if err != nil {
		return ast.Cast{}, err
	}

	if p.c.Kind != token.Comma {
		return ast.Cast{}, p.unexpected()
	}
	p.advance() // skip ","

	exp, err := p.Exp()
	if err != nil {
		return ast.Cast{}, err
	}

	if p.c.Kind != token.RightParen {
		return ast.Cast{}, p.unexpected()
	}
	p.advance() // skip ")"

	return ast.Cast{
		Type: typ,
		Exp:  exp,
	}, nil
}

func (p *Parser) Tint() (ast.Tint, diag.Error) {
	p.advance() // skip "tint"

	if p.c.Kind != token.LeftParen {
		return ast.Tint{}, p.unexpected()
	}
	p.advance() // skip "("

	typ, err := p.TypeSpec()
	if err != nil {
		return ast.Tint{}, err
	}

	if p.c.Kind != token.Comma {
		return ast.Tint{}, p.unexpected()
	}
	p.advance() // skip ","

	exp, err := p.Exp()
	if err != nil {
		return ast.Tint{}, err
	}

	if p.c.Kind != token.RightParen {
		return ast.Tint{}, p.unexpected()
	}
	p.advance() // skip ")"

	return ast.Tint{
		Type: typ,
		Exp:  exp,
	}, nil
}

func (p *Parser) Size() (ast.Size, diag.Error) {
	p.advance() // skip "#size"

	if p.c.Kind != token.LeftParen {
		return ast.Size{}, p.unexpected()
	}
	p.advance() // skip "("

	spec, err := p.TypeSpec()
	if err != nil {
		return ast.Size{}, err
	}

	if p.c.Kind != token.RightParen {
		return ast.Size{}, p.unexpected()
	}
	p.advance() // skip ")"

	return ast.Size{Exp: spec}, nil
}

func (p *Parser) List() (ast.List, diag.Error) {
	pin := p.c.Pin
	if p.c.Kind == token.Chunk {
		p.advance() // skip "[]"
		return ast.List{Pin: pin}, nil
	}

	if p.c.Kind != token.LeftSquare {
		return ast.List{}, p.unexpected()
	}
	p.advance() // skip "["

	var list []ast.Exp
	for {
		if p.c.Kind == token.RightSquare {
			p.advance() // skip ""
			return ast.List{
				Pin:  pin,
				Exps: list,
			}, nil
		}

		exp, err := p.Exp()
		if err != nil {
			return ast.List{}, err
		}
		list = append(list, exp)

		if p.c.Kind == token.Comma {
			p.advance() // skip ","
		} else if p.c.Kind == token.RightSquare {
			// will be skipped at next iteration
		} else {
			return ast.List{}, p.unexpected()
		}
	}
}
