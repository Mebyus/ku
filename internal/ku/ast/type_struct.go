package ast

import "github.com/mebyus/ku/internal/ku/sx"

// Struct represents struct type specifier.
//
// Formal definition:
//
//	Struct -> "struct" "{" { Field [ "," ] } "}"
type Struct struct {
	spec

	// Can be nil (if struct does not have fields).
	Fields []Field

	Pin sx.Pin
}

var _ TypeSpec = &Struct{}

// Field represents a single field in struct or union, or form type specifier.
//
// Formal definition:
//
//	Field  -> Name ":" TypeSpec
//	Name   -> word
type Field struct {
	Name string
	Type TypeSpec

	// pin of field name
	Pin sx.Pin
}
