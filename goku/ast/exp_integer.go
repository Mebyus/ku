package ast

import (
	"fmt"
	"strconv"

	"github.com/mebyus/ku/goku/enums/exk"
	"github.com/mebyus/ku/goku/source"
)

// IntKind indicates integer literal kind.
type IntKind uint32

const (
	// IntDec - decimal integer literal.
	IntDec IntKind = iota

	// IntHex - hexadecimal integer literal.
	IntHex

	// IntBin - binary integer literal.
	IntBin

	// IntOct - octal integer literal.
	IntOct
)

// Integer represents a single integer token usage inside the tree.
type Integer struct {
	// Integer value represented by token.
	Val uint64

	Pin source.Pin

	// Auxiliary information about the token.
	Aux uint32
}

var _ Exp = Integer{}

func (Integer) Kind() exk.Kind {
	return exk.Integer
}

func (n Integer) Span() source.Span {
	return source.Span{Pin: n.Pin, Len: uint32(len(n.String()))}
}

func (n Integer) String() string {
	k := n.IntKind()
	switch k {
	case IntDec:
		return strconv.FormatUint(n.Val, 10)
	case IntHex:
		return "0x" + strconv.FormatUint(n.Val, 16)
	case IntBin:
		return "0b" + strconv.FormatUint(n.Val, 2)
	case IntOct:
		return "0o" + strconv.FormatUint(n.Val, 8)
	default:
		panic(fmt.Sprintf("unexpected integer literal kind (=%d)", k))
	}
}

func (n Integer) IntKind() IntKind {
	return IntKind(n.Aux)
}
