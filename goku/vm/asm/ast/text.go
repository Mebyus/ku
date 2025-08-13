package ast

import "github.com/mebyus/ku/goku/compiler/srcmap"

type Text struct {
	Functions []Fun

	// Optional. Contains empty Name field if absent.
	Entry Entry
}

// Entry represents entrypoint construct in program text.
type Entry struct {
	Name string
	Pin  srcmap.Pin
}
