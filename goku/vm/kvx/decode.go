package kvx

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

func Load(path string) (*Program, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return Decode(file)
}

func Decode(in io.ReadSeeker) (*Program, error) {
	g := Decoder{in: in}
	var h Header

	err := g.header(&h)
	if err != nil {
		return nil, err
	}

	text, err := g.data(h.Text.Offset, h.Text.Size)
	if err != nil {
		return nil, err
	}

	data, err := g.data(h.Data.Offset, h.Data.Size)
	if err != nil {
		return nil, err
	}

	return &Program{
		Text: text,
		Data: data,

		GlobalSize: h.Global.Size,
	}, nil
}

type Decoder struct {
	buf []byte

	in io.ReadSeeker

	pos uint64
}

func (g *Decoder) header(h *Header) error {
	g.buf = make([]byte, HeaderSize)
	n, err := io.ReadFull(g.in, g.buf)
	g.pos += uint64(n)
	if err != nil {
		return err
	}

	magic := g.buf[:4]
	if string(magic) != Magic {
		return fmt.Errorf("bad magic: %X", magic)
	}

	version := val32(g.buf[4:8])
	if version != 0 {
		return fmt.Errorf("bad version: %d", version)
	}

	h.Text = SegmentHeader{
		Offset: val64(g.buf[8:16]),
		Size:   val32(g.buf[16:20]),
		Flags:  val32(g.buf[20:24]),
	}

	h.Data = SegmentHeader{
		Offset: val64(g.buf[24:32]),
		Size:   val32(g.buf[32:36]),
		Flags:  val32(g.buf[36:40]),
	}

	h.Global = SegmentHeader{
		Offset: val64(g.buf[40:48]),
		Size:   val32(g.buf[48:52]),
		Flags:  val32(g.buf[52:56]),
	}

	return nil
}

// allocates and reads exactly size bytes at specified offset
func (g *Decoder) data(offset uint64, size uint32) ([]byte, error) {
	if size == 0 {
		return nil, nil
	}

	_, err := g.in.Seek(int64(offset), io.SeekStart)
	if err != nil {
		return nil, err
	}
	g.pos = offset

	data := make([]byte, size)
	n, err := io.ReadFull(g.in, data)
	g.pos += uint64(n)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func val32(buf []byte) uint32 {
	return binary.LittleEndian.Uint32(buf)
}

func val64(buf []byte) uint64 {
	return binary.LittleEndian.Uint64(buf)
}
