package ast

import (
	"github.com/mebyus/ku/internal/fbuf"
	"github.com/mebyus/ku/internal/ku/sx"
)

// SymExp represents expression operand formed by usage of symbol name.
type SymExp struct {
	operand

	Name string

	Pin sx.Pin
}

// Explicit interface implementation check.
var _ Operand = &SymExp{}

func (s *SymExp) Format(buf *fbuf.Buffer) {
	buf.Puts(s.Name)
}

// SymAddr represents expression formed by taking a symbol's address.
type SymAddr struct {
	operand

	// symbol name which is being addressed
	Name string

	// name pin
	Pin sx.Pin
}

// Explicit interface implementation check.
var _ Operand = &SymAddr{}

func (s *SymAddr) Format(buf *fbuf.Buffer) {
	buf.Puts(s.Name)
	buf.Puts(".&")
}

// SymDeref represents expression formed by dereferencing a symbol.
type SymDeref struct {
	operand

	// symbol which is being dereferenced
	Name string

	// name pin
	Pin sx.Pin
}

// Explicit interface implementation check.
var _ Operand = &SymDeref{}

func (s *SymDeref) Format(buf *fbuf.Buffer) {
	buf.Puts(s.Name)
	buf.Puts(".*")
}

// SymCall represents expression operand formed by function-like
// call to a symbol.
type SymCall struct {
	operand

	// can be nil if call does not have arguments
	Args []Exp

	// name of the called symbol
	Name string

	// name pin
	Pin sx.Pin
}

// Explicit interface implementation check.
var _ Operand = &SymCall{}

func (s *SymCall) Format(buf *fbuf.Buffer) {
	buf.Puts(s.Name)
	if len(s.Args) == 0 {
		buf.Puts("()")
		return
	}

	buf.Putb('(')
	s.Args[0].Format(buf)
	for _, arg := range s.Args[1:] {
		buf.Puts(", ")
		arg.Format(buf)
	}
	buf.Putb(')')
}
