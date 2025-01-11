package lexer

import "github.com/mebyus/ku/goku/token"

func (lx *Lexer) Lex() token.Token {
	return token.Token{Kind: token.EOF, Pin: lx.pin()}
}
