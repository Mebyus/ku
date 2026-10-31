package sx

import (
	"hash/fnv"
	"io"
)

// Text represents something that contains source text.
// Most often text comes from a file, but sometimes it may be generated on the
// fly during compilation and not stored in a filesystem. Other cases include
// source text from strings which is helpful for automated tests.
type Text struct {
	// We use string to hold text (file) contents in order to avoid token data
	// allocations later. This way we can take string slices of original text
	// string as token data.
	Data string

	// Empty when text does not come from a file.
	Path string

	// Used for texts which do not originate from filesystem.
	//
	// For example text can be created from a string with
	// label attached as its "name" for debug output.
	label string

	// File extension with leading ".". Empty when text does not come from a file.
	//
	// Examples:
	//	- ".ku"
	//	- ".c"
	//	- ".h"
	Ext string

	// Contains consistent hash of bytes stored in Data field.
	Hash uint64

	// Assigned automatically when text is loaded by Pool.
	// Zero value is reserved for texts which are used for consistent testing.
	ID uint32
}

// Label returns source text label that should be used for presenting
// source position to the user.
//
// Uses either filesystem path or attached label if text was not loaded
// from file.
func (t *Text) Label() string {
	if t.Path != "" {
		return t.Path
	}
	return t.label
}

// hash provided string with the same method that is used
// for loading text from files
func hashText(text string) uint64 {
	h := fnv.New64a()
	io.WriteString(h, text)
	return h.Sum64()
}
