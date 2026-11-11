package ast

import (
	"fmt"

	"github.com/mebyus/ku/internal/fbuf"
	"github.com/mebyus/ku/internal/ku/sx"
)

// Charm (named after magical token or amulet) describes a special
// property of a symbol, struct or union field. Each valid charm is
// recognized by some subsystem inside compiler toolchain and can
// adjust compilation in some way. Hence the name "charm": it is
// meant to do something "magical" deep inside the compiler during
// program compilation.
//
// Syntax example:
//
//	#[export]             // make symbol visible in output object file
//	#[export("name")]     // same as above, but with specific linkname
//	#[entry]              // mark function as program entrypoint
//	#[entry(start)]       // same as above, but specify entrypoint style
//	#[link.name("other")] // specify internal linkname for function (mostly to circumvent cyclic imports restriction)
//	#[offset(0)]          // override default byte offset of struct field
//
// Each charm has a name which comes right after "#[" as one or more
// words separated by "." and zero or less arguments. Each argument
// can either be (in token terms) a string (ex. "abc"), word or integer.
type Charm struct {
	// Empty if charm has no args.
	Args []Arg

	// Name can be a single word or several words separated
	// by ".". Parser joins them together from individual tokens
	// if dots are present. Examples:
	//
	//	- export
	//	- extern
	//	- entry
	//	- link.name
	//	- link.lib
	//
	// Always not empty.
	Name string

	// Pin of first word in name.
	Pin sx.Pin
}

func (c *Charm) Format(buf *fbuf.Buffer) {
	buf.Puts("#[")
	buf.Puts(c.Name)
	if len(c.Args) != 0 {
		buf.Putb('(')
		c.Args[0].Format(buf)
		for _, a := range c.Args[1:] {
			buf.Puts(", ")
			a.Format(buf)
		}
		buf.Putb(')')
	}
	buf.Putb(']')
}

// Arg represents charm argument, which can be either:
//   - "string"
//   - word     (ex: main, weak)
//   - integer  (ex: 0, 1, 4)
type Arg struct {
	// String or word value.
	Text string

	// Integer value. Sadly go does not support union or sum types
	// to overlap this field with Text field in memory.
	Int uint64

	Pin sx.Pin

	// Discriminator for Arg value type.
	Type ArgType
}

func (a *Arg) Format(buf *fbuf.Buffer) {
	switch a.Type {
	case ArgWord:
		buf.Puts(a.Text)
	case ArgStr:
		buf.Putb('"')
		buf.Puts(a.Text)
		buf.Putb('"')
	case ArgInt:
		buf.Putn(a.Int)
	default:
		panic(fmt.Sprintf("invalid type (=%d)", a.Type))
	}
}

type ArgType uint32

// Arg types.
const (
	ArgWord ArgType = iota
	ArgStr
	ArgInt
)
