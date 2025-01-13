package origin

import (
	"fmt"
	"sort"
)

// Path a combination of import string and origin is called "unit path".
// It uniquely identifies a unit (local or foreign) within a project.
type Path struct {
	// Always empty if origin is empty.
	// Always not empty if origin is not empty.
	Import string

	Origin Origin
}

var Empty = Path{}

func (p Path) IsEmpty() bool {
	return p.Origin.IsEmpty()
}

func (p Path) String() string {
	return fmt.Sprintf("%s: %s", p.Origin, p.Import)
}

func Sort(p []Path) {
	sort.Slice(p, func(i, j int) bool {
		a := p[i]
		b := p[j]
		return Less(a, b)
	})
}

// Less returns true if a is less than b
func Less(a, b Path) bool {
	return a.Origin < b.Origin || (a.Origin == b.Origin && a.Import < b.Import)
}

func Local(s string) Path {
	return Path{
		Import: s,
		Origin: Loc,
	}
}

func Locals(ss []string) []Path {
	if len(ss) == 0 {
		return nil
	}

	paths := make([]Path, 0, len(ss))
	for _, s := range ss {
		paths = append(paths, Local(s))
	}
	return paths
}

type Set map[Path]struct{}

func NewSet() Set {
	return make(Set)
}

func (s Set) Add(p Path) {
	s[p] = struct{}{}
}

func (s Set) Has(p Path) bool {
	_, ok := s[p]
	return ok
}

func (s Set) Sorted() []Path {
	if len(s) == 0 {
		return nil
	}

	pp := make([]Path, 0, len(s))
	for p := range s {
		pp = append(pp, p)
	}
	Sort(pp)
	return pp
}
