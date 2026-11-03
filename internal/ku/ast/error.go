package ast

import (
	"github.com/mebyus/ku/internal/fbuf"
	"github.com/mebyus/ku/internal/ku/sx"
)

// InvExp represents invalid expression which is produced by error
// upon parsing part of expression.
//
// Can act as expression.
type InvExp struct {
	operand

	Pin sx.Pin
}

var _ Exp = &InvExp{}

func (v *InvExp) Format(buf *fbuf.Buffer) {
	buf.Puts("<|inv|>")
}

// InvType represents invalid (error upon parsing) type specifier.
//
// Can act as type specifier.
type InvType struct {
	spec

	Pin sx.Pin
}

var _ TypeSpec = &InvType{}

func (v *InvType) Format(buf *fbuf.Buffer) {
	buf.Puts("<|inv|>")
}
