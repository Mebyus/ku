package ast

import (
	"github.com/mebyus/ku/internal/ku/sx"
	"github.com/mebyus/ku/internal/ku/token"
)

type Error struct {
	// Tokens consumed during error recovery.
	Tokens []token.Token

	// Short error name/description to inform user what went wrong.
	Short string

	Pin sx.Pin
}

// ErrorExp represents error upon parsing part of expression.
//
// Can act as expression.
type ErrorExp struct {
	exp
	Error
}

var _ Exp = &ErrorExp{}

// ErrorType represents error upon parsing type specifier.
//
// Can act as type specifier.
type ErrorType struct {
	spec
	Error
}

var _ TypeSpec = &ErrorType{}
