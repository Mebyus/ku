package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/exk"
	"github.com/mebyus/ku/goku/compiler/source"
)

// Name represents a single word token usage inside the tree.
type Name struct {
	// String that constitues the word.
	Word string

	Pin source.Pin
}

var _ Exp = Name{}

func (Name) Kind() exk.Kind {
	return exk.Name
}

func (d Name) Span() source.Span {
	return source.Span{Pin: d.Pin, Len: uint32(len(d.Word))}
}

func (d Name) String() string {
	return d.Word
}
