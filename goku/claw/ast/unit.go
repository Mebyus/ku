package ast

import "github.com/mebyus/ku/goku/compiler/srcmap"

// Dir represents unit build directive.
type Dir interface {
	_dir()
}

// Embed this to quickly implement _dir() discriminator from Dir interface.
// Do not use it for anything else.
type nodeDir struct{}

func (nodeDir) _dir() {}

type Import struct {
	nodeDir

	// String literal value represented by token.
	Val string

	Pin srcmap.Pin
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

type Unit struct {
	Dirs []Dir
}
