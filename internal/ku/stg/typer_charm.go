package stg

import (
	"fmt"

	"github.com/mebyus/ku/internal/ku/ast"
	"github.com/mebyus/ku/internal/ku/enums/symk"
)

func (t *Typer) applyCharms(s *Symbol, charms []ast.Charm) {
	if len(charms) == 0 {
		return
	}
	t.resetCharms(len(charms))

	for _, c := range charms {
		_, ok := t.charms[c.Name]
		if ok {
			t.report(c.Pin, 0, fmt.Sprintf("duplicate \"%s\" charm on the same symbol", c.Name))
			continue
		}
		t.charms[c.Name] += 1

		switch c.Name {
		case "export":
			if len(c.Args) == 0 {
				t.report(c.Pin, 0, "expected export name as charm argument")
				continue
			}
			if len(c.Args) != 1 {
				t.report(c.Pin, 0, "too many arguments for \"%s\" charm", c.Name)
				continue
			}
			arg := c.Args[0]
			if arg.Type != ast.ArgStr {
				t.report(arg.Pin, 0, "expected string argument as export name")
				continue
			}
			if arg.Text == "" {
				t.report(arg.Pin, 0, "empty export name")
				continue
			}
			if s.Link != "" {
				t.report(c.Pin, 0, "link name was already specified as \"%s\"", s.Link)
				continue
			}

			s.Link = arg.Text
			s.Flags |= SymbolExport
			t.unit.Exports = append(t.unit.Exports, s)
		case "extern":
			panic("TODO")
		case "link.name":
			if len(c.Args) == 0 {
				t.report(c.Pin, 0, "expected link name as charm argument")
				continue
			}
			if len(c.Args) != 1 {
				t.report(c.Pin, 0, "too many arguments for \"%s\" charm", c.Name)
				continue
			}
			arg := c.Args[0]
			if arg.Type != ast.ArgStr {
				t.report(arg.Pin, 0, "expected string argument as link name")
				continue
			}
			if arg.Text == "" {
				t.report(arg.Pin, 0, "empty link name")
				continue
			}
			if s.Link != "" {
				t.report(c.Pin, 0, "link name was already specified as \"%s\"", s.Link)
				continue
			}

			s.Link = arg.Text
		case "entry":
			if len(c.Args) != 0 {
				panic("TODO")
			}

			if t.unit.Entry != nil {
				t.report(c.Pin, 0, "unit already has function \"%s\" marked as entrypoint", t.unit.Entry.Name)
				continue
			}
			if s.Kind != symk.Fun {
				t.report(c.Pin, 0, fmt.Sprintf("non-function symbol \"%s\" marked as entrypoint", s.Name))
				continue
			}

			t.unit.Entry = s
		default:
			t.report(c.Pin, 0, fmt.Sprintf("unknown \"%s\" charm", c.Name))
		}
	}
}

func (t *Typer) resetCharms(size int) {
	if t.charms == nil {
		t.charms = make(map[string]int, size)
	} else {
		clear(t.charms)
	}
}
