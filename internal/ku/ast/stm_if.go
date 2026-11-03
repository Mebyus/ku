package ast

import (
	"github.com/mebyus/ku/internal/fbuf"
	"github.com/mebyus/ku/internal/ku/sx"
)

// IfBlock represents "if-then" branching statement which
// does not have else block.
type IfBlock struct {
	stm

	// condition
	Exp Exp

	// then branch
	Body Block

	// "if" pin
	Pin sx.Pin
}

var _ Statement = &IfBlock{}

func (f *IfBlock) Format(buf *fbuf.Buffer) {
	buf.Puts("if ")
	f.Exp.Format(buf)
	buf.Space()
	f.Body.Format(buf)
}

// Branch represents "if-then-else" branching statement.
type Branch struct {
	stm

	// condition
	Exp Exp

	// true branch
	Body Block

	Pin sx.Pin

	// can be nil if statement does not have else branch
	Else *Block
}

var _ Statement = &Branch{}

func (f *Branch) Format(buf *fbuf.Buffer) {
	buf.Puts("if ")
	f.Exp.Format(buf)
	buf.Space()
	f.Body.Format(buf)
	buf.Puts(" else ")
	f.Else.Format(buf)
}

// IfFlow represents short (one line) form of "if" branching statement.
// It supports only restricted set of statements which affect control flow:
//
//   - return
//   - break
//   - gonext
//   - goto
type IfFlow struct {
	stm

	// condition
	Exp Exp

	// true branch
	Then Statement

	Pin sx.Pin
}

var _ Statement = &IfFlow{}

func (f *IfFlow) Format(buf *fbuf.Buffer) {
	buf.Puts("if ")
	f.Exp.Format(buf)
	buf.Space()
	f.Then.Format(buf)
}
