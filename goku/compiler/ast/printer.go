package ast

import (
	"io"
	"strconv"
)

type Printer struct {
	// Output buffer.
	buf []byte

	// Indentation buffer.
	//
	// Stores sequence of bytes which is used for indenting current line
	// in output. When a new line starts this buffer is used to add indentation.
	ib []byte
}

func (g *Printer) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(g.buf)
	return int64(n), err
}

func (g *Printer) Str() string {
	return string(g.buf)
}

// put decimal formatted integer into output buffer
func (g *Printer) putn(n uint64) {
	g.puts(strconv.FormatUint(n, 10))
}

// put string into output buffer
func (g *Printer) puts(s string) {
	g.buf = append(g.buf, s...)
}

// put single byte into output buffer
func (g *Printer) putb(b byte) {
	g.buf = append(g.buf, b)
}

func (g *Printer) put(b []byte) {
	g.buf = append(g.buf, b...)
}

func (g *Printer) nl() {
	g.putb('\n')
}

func (g *Printer) space() {
	g.putb(' ')
}

func (g *Printer) semi() {
	g.putb(';')
}

// increment indentation by one level.
func (g *Printer) inc() {
	g.ib = append(g.ib, '\t')
}

// decrement indentation by one level.
func (g *Printer) dec() {
	g.ib = g.ib[:len(g.ib)-1]
}

// add indentation to current line.
func (g *Printer) indent() {
	g.put(g.ib)
}
