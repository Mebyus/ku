package srcmap

import (
	"io"
	"path/filepath"
)

// Text represents something that contains source text.
// Most often text comes from a file, but sometimes it may be generated on the
// fly during compilation and not stored on a file system. Other cases include
// source text from strings which is helpful for automated tests.
type Text struct {
	Data []byte

	// Empty when text does not come from a file.
	Path string

	// File extension with leading ".". Empty when text does not come from a file.
	//
	// Examples:
	//	- ".ku"
	//	- ".c"
	//	- ".h"
	//	- ".kub"
	Ext string

	// Assigned automatically when text is loaded by Pool.
	// Zero value is reserved for texts which are used for consistent testing.
	ID uint32
}

// NewText constructs a new Text as though it comes from a file.
func NewText(path string, data []byte) *Text {
	return &Text{
		Data: data,
		Path: path,
		Ext:  filepath.Ext(path),
	}
}

func (t *Text) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(t.Data)
	return int64(n), err
}
