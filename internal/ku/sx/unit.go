package sx

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mebyus/ku/internal/ku/char"
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
	//
	// TODO: maybe just make it zero?
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

// ValidImport returns (s, "") if provided string was valid,
// original string is returned unchanged in this case.
//
// Returns (s1, <error text>) if provided string was invalid.
// Returned string s1 is a corrected import string. Which is
// a best guess what original string meant.
//
// Compiler should report produced error message and keep going
// when operating in error tolerant mode, but still must mark
// build as failed.
func ValidImport(s string) (string, string) {
	if s == "" {
		return "", "empty import string"
	}

	var errtext string

	n := len(s)
	s = strings.TrimSpace(s)
	if len(s) != n {
		errtext = "import string contains whitespace"
	}
	if s == "" {
		return "", errtext
	}

	// TODO: optimize to avoid allocations
	parts := strings.Split(s, "/")
	v := 0 // number of valid parts
	for i := range len(parts) {
		p := parts[i]
		if p == "" || strings.ReplaceAll(p, ".", "") == "" {
			// skip empty or containing only dots parts
			if errtext == "" {
				errtext = "import string contains empty parts"
			}
			continue
		}

		for i := range len(p) {
			c := p[i]
			if !char.IsAlphanum(c) {
				return "", "import string contains non-alphanumerical characters"
			}
		}

		parts[v] = p
		v += 1
	}

	if v == len(parts) && errtext == "" {
		return s, ""
	}

	s = strings.Join(parts[:v], "/")
	return s, errtext
}
