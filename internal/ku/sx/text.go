package sx

import (
	"hash/fnv"
	"io"

	"github.com/mebyus/ku/internal/ku/char"
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

	// Table with line offsets.
	//
	// Contains byte offsets at which line with specific index starts.
	//
	// Computed lazily, meaning that for newly loaded texts this field
	// will be nil and will only be filled when line information is needed.
	linetab []uint32

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

// LineTab returns table with line offsets. It stores start offset for each
// line by its zero-based index. For example:
//
//	tab := text.LineTab()
//	offset := tab[4] // gives starting offset of text line 5
//
// Returns nil only if text is empty.
func (t *Text) LineTab() []uint32 {
	if t.linetab != nil {
		return t.linetab
	}

	text := t.Data
	n := uint32(len(text))
	if n == 0 {
		return nil
	}

	tab := []uint32{0} // first line always starts at offset 0
	var i uint32
	for i < n {
		if text[i] == '\n' {
			// line starts at next byte after newline character
			tab = append(tab, i+1)
		}
		i += 1
	}

	t.linetab = tab
	return tab
}

// FindPos find text position (line + column) for the given
// byte offset.
func (t *Text) FindPos(p uint32) TextPos {
	text := t.Data
	tab := t.LineTab()
	line := FindLine(tab, p)

	i := tab[line]
	var column uint32 // current column, zero-based value
	for i < p {
		if char.IsCodePointStart(text[i]) {
			column += 1
		}
		i += 1
	}
	return TextPos{
		Line:   line,
		Column: column,
	}
}

// hash provided string with the same method that is used
// for loading text from files
func hashText(text string) uint64 {
	h := fnv.New64a()
	io.WriteString(h, text)
	return h.Sum64()
}
