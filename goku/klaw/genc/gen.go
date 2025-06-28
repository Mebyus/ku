package genc

import (
	"io"
	"strconv"
)

type Gen struct {
	// Output buffer.
	buf []byte

	// Indentation buffer.
	//
	// Stores sequence of bytes which is used for indenting current line
	// in output. When a new line starts this buffer is used to add indentation.
	ib []byte
}

func (g *Gen) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(g.Bytes())
	return int64(n), err
}

func (g *Gen) Output() string {
	return string(g.Bytes())
}

func (g *Gen) Bytes() []byte {
	return g.buf
}

func (g *Gen) empty() bool {
	return len(g.buf) == 0
}

// put decimal formatted integer into output buffer
func (g *Gen) putn(n uint64) {
	g.puts(strconv.FormatUint(n, 10))
}

func (g *Gen) putlen(l int) {
	g.putn(uint64(l))
}

// put string into output buffer
func (g *Gen) puts(s string) {
	g.buf = append(g.buf, s...)
}

// put single byte into output buffer
func (g *Gen) putb(b byte) {
	g.buf = append(g.buf, b)
}

func (g *Gen) put(b []byte) {
	g.buf = append(g.buf, b...)
}

func (g *Gen) nl() {
	g.putb('\n')
}

func (g *Gen) space() {
	g.putb(' ')
}

func (g *Gen) semi() {
	g.putb(';')
}

// increment indentation by one level.
func (g *Gen) inc() {
	g.ib = append(g.ib, '\t')
}

// decrement indentation by one level.
func (g *Gen) dec() {
	g.ib = g.ib[:len(g.ib)-1]
}

// add indentation to current line.
func (g *Gen) indent() {
	g.put(g.ib)
}
