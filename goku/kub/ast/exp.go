package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/bok"
	"github.com/mebyus/ku/goku/compiler/srcmap"
)

type Exp interface {
	_exp()
}

// Embed this to quickly implement _exp() discriminator from Exp interface.
// Do not use it for anything else.
type nodeExp struct{}

func (nodeExp) _exp() {}

type String struct {
	nodeExp

	// String literal value represented by token.
	Val string

	Pin srcmap.Pin
}

type Integer struct {
	nodeExp

	// Integer value represented by token.
	Val uint64

	Pin srcmap.Pin
}

type Word struct {
	// String that constitues the word.
	Str string

	Pin srcmap.Pin
}

// Name represents usage of something with a name.
// Each name consists of 1 or more parts, separated by period.
//
//	name
//	other.name
//	three.parts.name
type Name struct {
	nodeExp

	Parts []Word
}

type Binary struct {
	nodeExp

	Op BinOp

	// Left side of binary expression.
	A Exp

	// Right side of binary expression.
	B Exp
}

// BinOp represents binary operator inside expression.
type BinOp struct {
	Pin  srcmap.Pin
	Kind bok.Kind
}
