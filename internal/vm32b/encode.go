package vm32b

import (
	"io"
	"os"
	"slices"
)

// Encode a program into provided output writer in binary format.
//
// Format description:
//
//	KVX  - 3 bytes (magic signature)
//	TM   - 1 byte (target machine, 0 - vm32b)
//
//	NF   - 4 bytes (number of functions)
//	[]FE - NF * 16 bytes (array of finfo entries)
//
//	NI     - 4 bytes (number of instructions)
//	[]INST - rest (encoded instructions have varying length)
//
// Instructions which do not have an argument are encoded
// in single byte, which stores opcode:
//
//	0b0xxxxxxx - highest bit signals decoder that this instruction does not have argument
//
// Instructions with argument are encoded in 5 bytes, 1 byte for opcode
// and 4 for argument:
//
//	0b1xxxxxxx - highest bit signals decoder that this instruction have an argument
//
// All numbers, addresses and instruction arguments are stored in little
// endian format.
func Encode(out io.Writer, prog *Prog) error {
	var enc encoder
	enc.encode(prog)

	_, err := out.Write(enc.out)
	return err
}

// Save encode program into binary KVX format and save it
// into specified file.
func Save(prog *Prog, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return Encode(file, prog)
}

type encoder struct {
	out []byte
}

func (c *encoder) reset(size int) {
	c.out = slices.Grow(c.out[:0], size)
}

func (c *encoder) encode(prog *Prog) {
	// approximate encoded output size to reduce reallocations
	size := 3 + 1               // magic signature + target machine
	size += 4                   // number of functions
	size += 16 * len(prog.funs) // array of finfo entries
	size += 4                   // number of instructions
	size += 4 * len(prog.text)  // estimate for average encoded instruction size
	c.reset(size)

	c.magic(0)
	c.funs(prog.funs)
	c.text(prog.text)
}

func (c *encoder) funs(funs []finfo) {
	c.u32(uint32(len(funs)))

	for _, f := range funs {
		c.u32(uint32(f.addr))
		c.u32(f.narg)
		c.u32(f.nloc)
		c.u32(f.nret)
	}
}

func (c *encoder) text(text []inst) {
	c.u32(uint32(len(text)))

	for _, x := range text {
		arg := c.opcode(x.op)
		if arg {
			c.u32(x.arg)
		}
	}
}

// returns true if instruction with this opcode has an argument
func (c *encoder) opcode(op op) bool {
	if op < _opARG1 {
		c.out = append(c.out, byte(op))
		return false
	}

	c.out = append(c.out, byte((1<<7)|op))
	return true
}

func (c *encoder) magic(target byte) {
	c.out = append(c.out, 'K')
	c.out = append(c.out, 'V')
	c.out = append(c.out, 'X')
	c.out = append(c.out, target)
}

func (c *encoder) u32(v uint32) {
	c.out = append(c.out, byte(v))
	c.out = append(c.out, byte(v>>8))
	c.out = append(c.out, byte(v>>16))
	c.out = append(c.out, byte(v>>24))
}
