package ast

import (
	"github.com/mebyus/ku/goku/enums/exk"
	"github.com/mebyus/ku/goku/source"
)

// type NodeFamily uint32

// const (
// 	NodeExp NodeFamily = iota
// 	NodeStm
// 	Node
// )

type Node interface {
	// Family() NodeFamily

	Span() source.Span
	String() string
}

type Exp interface {
	Node

	Kind() exk.Kind
}
