package srcmap

import (
	"fmt"
	"strconv"
	"strings"
)

// PinMap is an abstraction over something that can translate Pin value to
// text position in file.
type PinMap interface {
	DecodePin(Pin) (FilePos, error)
}

type FilePos struct {
	Path string
	Pos  TextPos
}

func (p FilePos) String() string {
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

func ParseTextPos(s string) (TextPos, bool) {
	const maxUint32 = 0xFFFFFFFF
	lineString, columnString, ok := strings.Cut(s, ":")
	if !ok {
		return TextPos{}, false
	}
	line, err := strconv.ParseUint(lineString, 10, 64)
	if err != nil {
		return TextPos{}, false
	}
	column, err := strconv.ParseUint(columnString, 10, 64)
	if err != nil {
		return TextPos{}, false
	}

	line -= 1
	if line > maxUint32 {
		return TextPos{}, false
	}
	column -= 1
	if column > maxUint32 {
		return TextPos{}, false
	}

	return TextPos{Line: uint32(line), Column: uint32(column)}, true
}

func MustParseTextPos(s string) TextPos {
	pos, ok := ParseTextPos(s)
	if !ok {
		panic(fmt.Sprintf("must parse text position \"%s\"", s))
	}
	return pos
}

type Pos struct {
	// ID of Text.
	Text uint32

	// Byte offset into text data.
	Offset uint32
}

// PinTextMask returns a pin mask that should be used for pins that come from texts
// with specified id.
func PinTextMask(id uint32) Pin {
	return Pin(id) << 32
}

// Pin stores source position information in a compact way.
//
// Low 32-bits of this integer carry byte offset into text.
// High 32-bits carry text id.
type Pin uint64

func (p Pin) Pos() Pos {
	return Pos{
		Text:   uint32(p >> 32),
		Offset: uint32(p & 0xFFFFFFFF),
	}
}

// Span refers to continuous portion of source text.
type Span struct {
	Pin Pin
	Len uint32
}
