/*
Package fbuf implements general purpose text format buffer.
*/
package fbuf

import (
	"io"
	"slices"
	"strconv"
)

// Buffer implements many primitive convenience methods for
// formatting source text from tree-like structures.
type Buffer struct {
	// output buffer with raw bytes
	out []byte

	// indentation level
	il int
}

var _ io.Writer = &Buffer{}
var _ io.WriterTo = &Buffer{}

// Reset drop stored output, but keep underlying memory.
// If necessary preallocate enough memory to store next
// specified number of bytes.
func (g *Buffer) Reset(size int) {
	g.out = slices.Grow(g.out[:0], size)
}

// Write implements io.Writer interface.
func (g *Buffer) Write(data []byte) (int, error) {
	g.out = append(g.out, data...)
	return len(data), nil
}

// WriteTo implements io.WriterTo interface.
func (g *Buffer) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(g.out)
	return int64(n), err
}

// Bytes returns stored output for reading. Callers must not take
// ownership of returned slice.
func (g *Buffer) Bytes() []byte {
	if len(g.out) == 0 {
		return nil
	}
	return g.out
}

// Putb append single into stored output.
func (g *Buffer) Putb(b byte) {
	g.out = append(g.out, b)
}

// Puts append string into stored output.
func (g *Buffer) Puts(s string) {
	g.out = append(g.out, s...)
}

// Putn append decimal formatted integer into stored output.
func (g *Buffer) Putn(n uint64) {
	g.Puts(strconv.FormatUint(n, 10))
}

// Space append space character into stored output.
func (g *Buffer) Space() {
	g.Putb(' ')
}

// Semi append semicolon character into stored output.
func (g *Buffer) Semi() {
	g.Putb(';')
}

// Line append newline character (\n) into stored output.
func (g *Buffer) Line() {
	g.Putb('\n')
}

// Indent append indentation (with current level) into stored output.
func (g *Buffer) Indent() {
	for range g.il {
		g.Putb('\t')
	}
}

// Inc increase indentation level by 1.
func (g *Buffer) Inc() {
	g.il += 1
}

// Dec decrease indentation level by 1.
func (g *Buffer) Dec() {
	g.il -= 1
}
