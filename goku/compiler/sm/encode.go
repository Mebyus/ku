package sm

import (
	"fmt"
	"io"
	"slices"
)

// Encode line map of every text inside a pool.
func Encode(w io.Writer, p *Pool) error {
	var g Encoder
	g.encode(p.list)

	_, err := w.Write(g.header)
	if err != nil {
		return err
	}
	_, err = w.Write(g.names)
	if err != nil {
		return err
	}
	_, err = w.Write(g.lm)
	if err != nil {
		return err
	}
	return nil
}

type Encoder struct {
	// encoded header
	header []byte

	// concatenated file path strings data
	names []byte

	// line map
	lm []byte

	// total length of all path strings
	plen int

	// total size of all text data
	total int

	// current line map encoding position
	p int
}

// list must be ordered by Text.ID
func (g *Encoder) encode(list []*Text) {
	const debug = true

	g.calc(list)

	g.encodeNames(list)
	g.encodeTextMaps(list)

	if debug {
		fmt.Printf("src total size: %d\n", g.total)
		fmt.Printf("lm  total size: %d\n", len(g.lm))
		fmt.Printf("avg encode density: %d\n", g.total/len(g.lm))
	}
}

// calculate various sizes which are needed later during encoding
func (g *Encoder) calc(list []*Text) {
	plen := 0
	total := 0
	for _, t := range list {
		plen += len(t.Path)
		total += len(t.Data)
	}

	g.plen = plen
	g.total = total
}

func (g *Encoder) encodeNames(list []*Text) {
	g.names = slices.Grow(g.names, g.plen)

	for _, t := range list {
		g.names = append(g.names, t.Path...)
	}
}

func (g *Encoder) encodeTextMaps(list []*Text) {
	// prealloc space for line map
	g.lm = slices.Grow(g.lm, g.total/20) // 20 is an empirical constant for average encode density

	for _, t := range list {
		g.encodeTextMap(t)
	}
}

func (g *Encoder) encodeTextMap(t *Text) {
	d := t.Data
	i := 0
	n := 0    // number of bytes in current line
	prev := 0 // previous line start
	for i < len(d) {
		c := d[i]

		if c == '\n' {
			n = i - prev
			prev = i + 1
			g.putVarUint32IntoTextMap(uint32(n))
		}
		i += 1
	}
}

func (g *Encoder) putVarUint32IntoTextMap(v uint32) {
	if cap(g.lm)-len(g.lm) < 5 {
		g.lm = append(g.lm, 0, 0, 0, 0, 0) // 5 * 7 bits = 35 bits is enough to encode 32 bit integer
	} else {
		g.lm = g.lm[:cap(g.lm)]
	}
	n := putVarUint32(g.lm[g.p:], v)
	g.p += n
	g.lm = g.lm[:g.p]
}

// returns number of bytes used
func putVarUint32(b []byte, v uint32) int {
	i := 0
	for v >= 0x80 {
		b[i] = 0x80 | uint8(v&0x7F)
		v >>= 7
		i += 1
	}
	b[i] = uint8(v)
	return i + 1
}

// returns decoded value and number of bytes read
func getVarUint32(b []byte) (uint32, int) {
	var v uint32
	i := 0
	for {
		c := b[i]
		v |= uint32(c&0x7F) << (7 * i)
		i += 1

		if c&0x80 == 0 {
			break
		}
	}
	return v, i
}
