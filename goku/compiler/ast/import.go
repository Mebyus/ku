package ast

import (
	"github.com/mebyus/ku/goku/compiler/srcmap"
	"github.com/mebyus/ku/goku/compiler/srcmap/origin"
)

// <ImportBlock> = "import" [ <ImportOrigin> ] "(" { <ImportSpec> } ")"
//
// <ImportOrigin> = "std" | "pkg" | "loc"
//
// If <ImportOrigin> is absent in block, then it is equivalent to <ImportOrigin> = "loc".
// Canonical ku code style omits import origin in such cases, instead of specifying
// it to "loc" explicitly.
type ImportBlock struct {
	Imports []Import

	Origin origin.Origin
}

// <Import> = <Name> "=>" <ImportString>
//
// <ImportString> = <String> (cannot be empty)
type Import struct {
	Name   Word
	String ImportString
}

type ImportString struct {
	Pin srcmap.Pin
	Str string
}
