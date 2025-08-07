package kvx

import (
	"encoding/binary"
	"io"
	"os"
)

func Save(path string, prog *Program) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return Encode(file, prog)
}

func Encode(out io.Writer, prog *Program) error {
	file := NewFile(prog)
	g := Encoder{out: out}

	g.bufString(Magic)
	g.bufVal32(file.Header.Version)

	g.bufSegmentHeader(file.Header.Text)
	g.bufSegmentHeader(file.Header.Data)
	g.bufSegmentHeader(file.Header.Global)
	err := g.flush()
	if err != nil {
		return err
	}

	err = g.writeAt(prog.Text, file.Header.Text.Offset)
	if err != nil {
		return err
	}
	err = g.writeAt(prog.Data, file.Header.Data.Offset)
	if err != nil {
		return err
	}
	return nil
}

type Encoder struct {
	buf []byte

	out io.Writer

	// Total number of bytes written to out.
	pos uint64
}

func (g *Encoder) bufSegmentHeader(h SegmentHeader) {
	g.bufVal64(h.Offset)
	g.bufVal32(h.Size)
	g.bufVal32(h.Flags)
}

func (g *Encoder) write(data []byte) error {
	n, err := g.out.Write(data)
	g.pos += uint64(n)
	return err
}

func (g *Encoder) flush() error {
	err := g.write(g.buf)
	if err != nil {
		return err
	}
	g.buf = g.buf[:0]
	return nil
}

// Add padding to out until offset is reached and then write
// data to out.
//
// Panics if specified offset is less than number of bytes already
// written to out.
func (g *Encoder) writeAt(data []byte, offset uint64) error {
	if offset < g.pos {
		panic("invalid offset")
	}

	if offset > g.pos {
		n := offset - g.pos
		for range n {
			g.buf = append(g.buf, 0)
		}
		err := g.flush()
		if err != nil {
			return err
		}
	}

	return g.write(data)
}

func (g *Encoder) bufString(s string) {
	g.buf = append(g.buf, s...)
}

func (g *Encoder) bufVal32(v uint32) {
	g.buf = binary.LittleEndian.AppendUint32(g.buf, v)
}

func (g *Encoder) bufVal64(v uint64) {
	g.buf = binary.LittleEndian.AppendUint64(g.buf, v)
}
