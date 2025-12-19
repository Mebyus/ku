package sm

import (
	"fmt"
	"io"

	"github.com/mebyus/ku/goku/compiler/char"
)

type RenderParams struct {
	// Inserted before each line in formatted output.
	Prefix string

	Window WindowParams

	Text *Text
}

// TargetLineWindow combines a sequence of text lines around target position.
type TargetLineWindow struct {
	// Each line does not contain newline character at the end.
	Lines [][]byte

	// Index of target line inside Lines slice.
	Target uint32

	// Line number of the Lines[0] element.
	// Zero-based numeration.
	Start uint32
}

func Render(w io.Writer, params RenderParams) error {
	return nil
}

type WindowParams struct {
	Offset uint32

	MaxLinesBefore uint32
	MaxLinesAfter  uint32
}

// FindTargetLineWindow returns a window into text with several lines before
// and after line containing specified offset.
//
// Panics if offset is outside of text slice.
func FindTargetLineWindow(text []byte, params WindowParams) TargetLineWindow {
	targetLine := FindLineNumberAtOffset(text, params.Offset)
	_ = targetLine
	return TargetLineWindow{}
}

// FindLineNumberAtOffset returns zero-based line number of position in text.
// Position is specified by its offset.
//
// Panics if offset is outside of text.
func FindLineNumberAtOffset(text []byte, offset uint32) uint32 {
	var i uint32
	var line uint32 // current line, zero-based value
	for i < offset {
		if text[i] == '\n' {
			line += 1
		}
		i += 1
	}
	return line
}

// FindTextPos returns (line, column) text position for specified offset into text.
//
// Panics if offset is outside of text.
func FindTextPos(text []byte, offset uint32) TextPos {
	var i uint32
	var line uint32  // current line, zero-based value
	var start uint32 // current line start offset
	for i < offset {
		if text[i] == '\n' {
			start = i + 1
			line += 1
		}
		i += 1
	}

	i = start
	var column uint32 // current column, zero-based value
	for i < offset {
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

// FindLineOffset returns offset of line start before the given position.
//
// Panics if offset is outside of text.
func FindLineOffset(text []byte, offset uint32) uint32 {
	i := offset
	if i == uint32(len(text)) {
		if i == 0 {
			return 0
		}
		if text[i-1] == '\n' {
			return i
		}
	}
	for i != 0 {
		i -= 1
		if text[i] == '\n' {
			break
		}
	}
	if i == 0 {
		return 0
	}
	return i + 1
}

// FindNextLineOffset returns offset of next line start after the given position.
//
// Panics if offset is outside of text.
func FindNextLineOffset(text []byte, offset uint32) uint32 {
	return 0
}

// FindOffset tries to find byte offset of position (line, column) inside the text.
//
// If specified position exists within the text then (offset, true) is returned.
// Otherwise function returns (0, false).
func FindOffset(text []byte, pos TextPos) (uint32, bool) {
	var i uint32
	var line uint32 // current line, zero-based value
	for i < uint32(len(text)) && line < pos.Line {
		if text[i] == '\n' {
			line += 1
		}
		i += 1
	}
	if line < pos.Line {
		return 0, false
	}

	var col uint32 // current column, zero-based value
	for i < uint32(len(text)) && col < pos.Column {
		if text[i] == '\n' {
			return 0, false
		}
		if char.IsCodePointStart(text[i]) {
			col += 1
		}
		i += 1
	}
	if col < pos.Column {
		return 0, false
	}
	for i < uint32(len(text)) && char.InsideCodePoint(text[i]) {
		i += 1
	}
	return i, true
}

func MustFindOffset(text []byte, pos TextPos) uint32 {
	offset, ok := FindOffset(text, pos)
	if !ok {
		panic(fmt.Sprintf("must find offset of position \"%s\" inside text (len=%d)", pos, len(text)))
	}
	return offset
}
