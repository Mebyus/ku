package sm

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/mebyus/ku/goku/compiler/char"
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
type Origin uint8

const (
	// Zero value of Origin. Valid only as intermidiate value (e.g. check if struct was
	// filled or just created with zero values). Typically should not be used directly.
	empty Origin = iota

	// Std denotes import context of standard library.
	Std

	// Pkg denotes import context of units from other (third-party) projects.
	// Can be managed via "pkg.ku" file in project's root directory.
	Pkg

	// Loc denotes import context of units "local" to current project.
	// In contrast with other origins it is specified by omitting origin name.
	// in import block:
	//
	//	import {
	//		// import entries
	//	}
	Loc

	// Main denotes unit path of main unit. Main units cannot be imported by
	// other units.
	Main
)

func (o Origin) IsEmpty() bool {
	return o == empty
}

var text = [...]string{
	empty: "<nil>",

	Std:  "std",
	Pkg:  "pkg",
	Loc:  "loc",
	Main: "main",
}

func (o Origin) String() string {
	return text[o]
}

// Parse attempts to parse origin from a string. Second return value indicates
// whether this operation was successful or not
//
// Empty string yields local origin
func ParseOrigin(s string) (Origin, bool) {
	switch s {
	case "", "loc":
		return Loc, true
	case "std":
		return Std, true
	case "pkg":
		return Pkg, true
	default:
		return empty, false
	}
}

// UnitPath a combination of import string and origin is called "unit path".
// It uniquely identifies a unit (local or foreign) within a project.
type UnitPath struct {
	// Always empty if origin is empty.
	// Always not empty if origin is not empty.
	Import string

	Origin Origin
}

func (p *UnitPath) IsEmpty() bool {
	return p.Origin.IsEmpty()
}

func (p *UnitPath) String() string {
	return fmt.Sprintf("<%s> %s", p.Origin, p.Import)
}

func Sort(p []UnitPath) {
	sort.Slice(p, func(i, j int) bool {
		a := p[i]
		b := p[j]
		return Less(a, b)
	})
}

// Less returns true if a is less than b
func Less(a, b UnitPath) bool {
	return a.Origin < b.Origin || (a.Origin == b.Origin && a.Import < b.Import)
}

func Local(s string) UnitPath {
	return UnitPath{
		Import: s,
		Origin: Loc,
	}
}

func Locals(ss []string) []UnitPath {
	if len(ss) == 0 {
		return nil
	}

	paths := make([]UnitPath, 0, len(ss))
	for _, s := range ss {
		paths = append(paths, Local(s))
	}
	return paths
}

func CheckImportString(s string) error {
	if s == "" {
		return errors.New("empty import string")
	}

	split := strings.Split(s, "/")
	if len(split) == 0 {
		panic("impossible for non-empty string")
	}

	for _, part := range split {
		if part == "" {
			return errors.New("import string contains empty part")
		}
		if strings.Contains(part, "..") {
			return errors.New("import string part contains \"..\"")
		}
		if !char.IsLatinLetter(part[0]) {
			return fmt.Errorf("import string part starts from '%c' character", s[0])
		}

		for i := range len(part) {
			c := part[i]
			if char.IsAlphanum(c) || c == '.' || c == '-' {
				continue
			}

			return fmt.Errorf("import string part contains '%c' character", s[0])
		}
	}

	return nil
}

type PathSet map[UnitPath]struct{}

func NewPathSet() PathSet {
	return make(PathSet)
}

func (s PathSet) Add(p UnitPath) {
	s[p] = struct{}{}
}

func (s PathSet) Has(p UnitPath) bool {
	_, ok := s[p]
	return ok
}

func (s PathSet) Clear() {
	clear(s)
}

func (s PathSet) Sorted() []UnitPath {
	if len(s) == 0 {
		return nil
	}

	pp := make([]UnitPath, 0, len(s))
	for p := range s {
		pp = append(pp, p)
	}
	Sort(pp)
	return pp
}
