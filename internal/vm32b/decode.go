package vm32b

import "errors"

// DecodeFromBytes decodes program from encoded KVX format.
func DecodeFromBytes(data []byte) (*Prog, error) {
	dec := decoder{data: data}
	err := dec.decode()
	if err != nil {
		return nil, err
	}

	return &dec.prog, nil
}

type decoder struct {
	prog Prog

	data []byte

	// current decoder position
	p int
}

func (c *decoder) decode() error {
	if c.left() <= 0 {
		return errors.New("empty input")
	}

	// minimum size of KVX encoded program
	const minSize = 3 + 1 + // magic + target machine
		4 + // number of functions
		4 // number of instructions

	if c.left() < minSize {
		return errors.New("unrecognized format")
	}

	if !c.verify("KVX\x00") {
		return errors.New("unrecognized format")
	}

	nfun := c.u32()
	if nfun == 0 {
		return errors.New("empty finfo section")
	}

	const finfoSize = 16
	if c.left() < finfoSize*int(nfun) {
		return errors.New("malformed finfo section")
	}

	funs := make([]finfo, 0, nfun)
	for range nfun {
		funs = append(funs, c.finfo())
	}

	if c.left() < 4 {
		return errors.New("malformed text section")
	}

	n := c.u32()           // number of instructions in text
	if c.left() < int(n) { // minimum number of bytes for n instructions
		return errors.New("malformed text section")
	}

	text := make([]inst, 0, n)
	for range n {
		x, err := c.inst()
		if err != nil {
			return err
		}
		text = append(text, x)
	}

	c.prog.funs = funs
	c.prog.text = text
	return nil
}

func (c *decoder) finfo() finfo {
	addr := c.u32()
	narg := c.u32()
	nloc := c.u32()
	nret := c.u32()

	return finfo{
		addr: ia(addr),
		narg: narg,
		nloc: nloc,
		nret: nret,
	}
}

func (c *decoder) inst() (inst, error) {
	if c.left() < 1 {
		return inst{}, errors.New("malformed instruction")
	}

	b := c.data[c.p]
	c.p += 1

	if b&(1<<7) == 0 {
		return inst{op: op(b)}, nil
	}

	if c.left() < 4 {
		return inst{}, errors.New("malformed instruction")
	}
	arg := c.u32()

	return inst{
		op:  op(b & 0b01111111),
		arg: arg,
	}, nil
}

// consume number of bytes equal to specified string byte length
// and check that consumed data is equal to the string
//
// does no boundary checks
func (c *decoder) verify(s string) bool {
	for k := range len(s) {
		i := c.p + k
		if c.data[i] != s[k] {
			return false
		}
	}

	c.p += len(s)
	return true
}

func (c *decoder) u32() uint32 {
	b := c.data[c.p:]
	_ = b[3]

	v := uint32(b[0]) | (uint32(b[1]) << 8) | (uint32(b[2]) << 16) | (uint32(b[3]) << 24)

	c.p += 4
	return v
}

// returns number of undecoded bytes left
func (c *decoder) left() int {
	return len(c.data) - c.p
}
