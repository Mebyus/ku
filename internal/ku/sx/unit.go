package sx

import (
	"errors"
	"fmt"
	"strings"
)

// Origin denotes context in which import string should be resolved.
// Resolution mechanisms and rules may significantly vary between
// different origins.
//
// Origin is specified at the start of import block, like this:
//
//	import std {
//		// import entries
//	}
//
// Zero value of Origin indicates that it was filled incorrectly.
// Typically should not be used directly.
type Origin uint8

const (
	// Std denotes import context of standard library.
	Std Origin = 1 + iota

	// Pkg denotes import context of units from other (third-party) projects.
	Pkg

	// Loc denotes import context of units "local" to current project.
	// In contrast with other origins it is specified by omitting origin name.
	// in import block:
	//
	//	import {
	//		// import entries
	//	}
	Loc

	// Bad (invalid) origin.
	Bad
)

var originText = [...]string{
	Std: "std",
	Pkg: "pkg",
	Loc: "loc",
	Bad: "bad",
}

func (o Origin) String() string {
	err := o.Valid()
	if err != nil {
		return "???"
	}
	return originText[o]
}

func (o Origin) Valid() error {
	if o == 0 {
		return errors.New("empty origin")
	}
	if o > Loc {
		return fmt.Errorf("unknown origin (=%d)", o)
	}
	return nil
}

// Path combination of import string and origin is called "unit path".
// It uniquely identifies a unit (local or foreign) within a project.
//
// This implementation stores origin in the first byte of underlying string.
// Import string comes directly after that. It is designed for ease of use
// with native Go maps and comparison.
type Path string

func MakePath(o Origin, s string) Path {
	var g strings.Builder
	g.Grow(1 + len(s))

	_ = g.WriteByte(byte(o))
	_, _ = g.WriteString(s)
	return Path(g.String())
}

func (p Path) String() string {
	if p == "" {
		return "???"
	}

	o, s := p.Import()
	if s == "" {
		if o == Loc {
			return "???"
		}
		return o.String() + ": ???"
	}

	if o == Loc {
		return s
	}
	return fmt.Sprintf("%s: %s", o, s)
}

// Import extracts origin and import string from unit path.
func (p Path) Import() (Origin, string) {
	if p == "" {
		return 0, ""
	}
	return Origin(p[0]), string(p[1:])
}

// Name returns last part (after final "/") of import string.
// It can act as unit short name for the purpose of aliasing.
func (p Path) Name() string {
	s := string(p[1:])
	i := strings.LastIndexByte(s, '/')
	if i < 0 {
		return s
	}
	return s[i+1:]
}
