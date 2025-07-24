package ast

import (
	"github.com/mebyus/ku/goku/compiler/srcmap"
	"github.com/mebyus/ku/goku/compiler/srcmap/origin"
)

// Dir represents unit build directive.
type Dir interface {
	_dir()
}

// Embed this to quickly implement _dir() discriminator from Dir interface.
// Do not use it for anything else.
type nodeDir struct{}

func (nodeDir) _dir() {}

// Import represents single import (without block).
type Import struct {
	nodeDir

	// String literal value represented by token.
	Val string

	Pin srcmap.Pin

	Origin origin.Origin
}

type ImportBlock struct {
	nodeDir

	Imports []ImportString

	Origin origin.Origin
}

type ImportString struct {
	Pin srcmap.Pin
	Val string
}

type Include struct {
	nodeDir

	// String literal value represented by token.
	Val string

	Pin srcmap.Pin
}

type Block struct {
	Dirs []Dir

	Pin srcmap.Pin
}

type Test struct {
	nodeDir

	Block
}

type Exe struct {
	nodeDir

	Block
}

type Unit struct {
	Dirs []Dir
}
