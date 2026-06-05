package sx

import (
	"strconv"

	"github.com/mebyus/ku/internal/ku/char"
)

// PinTextMask returns a pin mask that should be used for pins that come from texts
// with specified id.
func PinTextMask(id uint32) Pin {
	return Pin(id) << 32
}

// Pin is encoded form of Pos.
// Pin stores source position information in a compact way.
//
// Low 32-bits of this integer carry byte offset into text.
// High 32-bits carry text id.
//
// Zero value is reserved for things that do not have position
// in source (ex: generated symbols and nodes).
type Pin uint64

// Pos is decoded form of Pin.
// Pos stores source position as text id + byte offset.
type Pos struct {
	// Text id.
	Text uint32

	// Byte offset into text data.
	Offset uint32
}

func (p Pin) Pos() Pos {
	return Pos{
		Text:   uint32(p >> 32),
		Offset: uint32(p & 0xFFFFFFFF),
	}
}

// FilePos refers to position inside text file as path + (line, column).
type FilePos struct {
	Path string
	Pos  TextPos
}

func (p FilePos) String() string {
	if p.Path == "" {
		return "???:" + p.Pos.String()
	}
	return p.Path + ":" + p.Pos.String()
}

// TextPos refers to position inside a text as (line, column) pair.
type TextPos struct {
	// Contains zero-based value.
	Line uint32

	// Contains zero-based value.
	Column uint32
}

func (p TextPos) String() string {
	return strconv.FormatUint(uint64(p.Line)+1, 10) + ":" + strconv.FormatUint(uint64(p.Column)+1, 10)
}

// FindTextPos returns (line, column) text position for specified offset into text.
//
// Panics if offset is outside of text.
func FindTextPos(text string, offset uint32) TextPos {
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
