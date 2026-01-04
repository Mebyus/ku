package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/stk"
	"github.com/mebyus/ku/goku/compiler/sm"
	"github.com/mebyus/ku/goku/compiler/token"
)

type Asm struct {
	Nodes []AsmNode

	Inits []RegInit
	Outs  []RegOut

	Clobbers []Clobber

	Arch Word
}

// RegInit register init expression in asm block.
type RegInit struct {
	// Register name.
	Name Word

	// Initial register value is taken from this symbol.
	//
	// We intentionally allow only symbols here and not arbitrary epressions
	// to keep things simple.
	Symbol Word
}

// RegOut register output expression in asm block.
type RegOut struct {
	// Register name.
	Name Word

	// Output symbol.
	Symbol Word
}

// Clobber register or memory clobber in asm block.
type Clobber struct {
	// Register name or memory.
	Name Word
}

// Explicit interface implementation check.
var _ Statement = &Asm{}

func (*Asm) Kind() stk.Kind {
	return stk.Asm
}

func (a *Asm) Span() sm.Span {
	return sm.Span{Pin: a.Arch.Pin}
}

func (a *Asm) String() string {
	panic("not implemented")
	// var g Printer
	// g.Assign(a)
	// return g.Output()
}

// AsmNode instruction or label placement.
type AsmNode interface{}

// AsmPlace label placement in asm block.
type AsmPlace struct {
	Name Word
}

// AsmInst instruction in asm block.
type AsmInst struct {
	// Tokens after mnemonic until ";".
	Rest []token.Token

	Mnemonic string

	Pin sm.Pin
}
