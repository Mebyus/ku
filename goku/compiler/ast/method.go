package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/tnk"
	"github.com/mebyus/ku/goku/compiler/srcmap"
)

// Formal definition:
//
//	Method => "fun" "(" Receiver ")" Name Signature Body
type Method struct {
	Signature Signature

	Body Block
	Name Word

	Receiver Receiver

	Traits
}

var _ Top = Method{}

func (Method) Kind() tnk.Kind {
	return tnk.Method
}

func (m Method) Span() srcmap.Span {
	return m.Name.Span()
}

func (m Method) String() string {
	var g Printer
	g.Method(m)
	return g.Output()
}

type ReceiverKind uint8

const (
	ReceiverVal ReceiverKind = 1 + iota
	ReceiverRef
	ReceiverPtr
)

// Formal definition:
//
//	Receiver -> [ Name ] [ "*" | "&" ] TypeName
type Receiver struct {
	Name Word
	Kind ReceiverKind
}
