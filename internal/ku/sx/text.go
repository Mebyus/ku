package sx

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
