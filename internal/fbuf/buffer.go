/*
Package fbuf implements general purpose text format buffer.
*/
package fbuf

import (
	"io"
	"os"
	"slices"
	"strconv"
)

// Buffer implements many primitive convenience methods for
// formatting source text from tree-like structures.
//
// Zero value is ready-to-use, empty buffer.
type Buffer struct {
	// output buffer with raw bytes
	out []byte

	// indentation level
	il int
}

var _ io.Writer = &Buffer{}
var _ io.WriterTo = &Buffer{}
var _ io.ReaderFrom = &Buffer{}

// Reset drop stored output, but keep underlying memory.
// If necessary preallocate enough memory to store next
// specified number of bytes.
func (g *Buffer) Reset(size int) {
	g.il = 0
	g.out = slices.Grow(g.out[:0], size)
}

// Grow (preallocate) more space in internal buffer for
// stored output by at least specified number of bytes.
//
// Does not affect already stored bytes or result of
// future writes/appends.
func (g *Buffer) Grow(size int) {
	g.out = slices.Grow(g.out, size)
}

// Write implements io.Writer interface.
func (g *Buffer) Write(data []byte) (int, error) {
	g.out = append(g.out, data...)
	return len(data), nil
}

// WriteTo implements io.WriterTo interface.
func (g *Buffer) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(g.Bytes())
	return int64(n), err
}

// ReadFrom implements io.ReaderFrom interface.
func (g *Buffer) ReadFrom(r io.Reader) (int64, error) {
	var buf [1 << 14]byte

	var n int64
	for {
		k, err := r.Read(buf[:])
		n += int64(k)
		g.Write(buf[:k])

		if err != nil {
			if err == io.EOF {
				return n, nil
			}
			return n, err
		}
	}
}

// ReadFile reads data from specified file and appends it
// into stored output. This method is very similar to ReadFrom
// but operates with file instead of io.Reader.
//
// Returns number of bytes read from file.
func (g *Buffer) ReadFile(path string) (int64, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return 0, err
	}
	size := stat.Size()
	if size <= 0 {
		return 0, nil
	}

	g.Grow(int(size))
	return g.ReadFrom(file)
}

// Len returns number of bytes currently stored in output buffer.
func (g *Buffer) Len() int {
	return len(g.out)
}

// Bytes returns stored output for reading. Callers must not take
// ownership of returned slice.
func (g *Buffer) Bytes() []byte {
	if len(g.out) == 0 {
		return nil
	}
	return g.out
}

// Putb append single byte into stored output.
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

// Put append raw bytes into stored output.
func (g *Buffer) Put(data []byte) {
	_, _ = g.Write(data)
}

// RepPutb append single byte into stored output specified number of times.
func (g *Buffer) RepPutb(b byte, n int) {
	if n <= 0 {
		return
	}
	g.Grow(n)

	for range n {
		g.Putb(b)
	}
}

// Puts append string into stored output specified number of times.
func (g *Buffer) RepPuts(s string, n int) {
	if n <= 0 || s == "" {
		return
	}
	g.Grow(n * len(s))

	for range n {
		g.Puts(s)
	}
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
	g.RepPutb('\t', g.il)
}

// RepSpace append space character into stored output specified number of times.
func (g *Buffer) RepSpace(n int) {
	g.RepPutb(' ', n)
}

// Inc increase indentation level by 1.
func (g *Buffer) Inc() {
	g.il += 1
}

// Dec decrease indentation level by 1.
func (g *Buffer) Dec() {
	g.il -= 1
}
