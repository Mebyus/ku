package parse

import (
	"errors"
	"fmt"

	"github.com/mebyus/ku/goku/butler"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/sm"
	"github.com/mebyus/ku/goku/kub/ast"
	"github.com/mebyus/ku/goku/kub/parser"
)

var paramType = butler.Param{
	Name:    "type",
	Alias:   "t",
	Desc:    "Specifies type of file to be parsed",
	Default: "package",
	Kind:    butler.String,
}

var Butler = &butler.Butler{
	Name: "parse",

	Short: "parse a given source file and print simplified AST of package",
	Usage: "[options] [file]",

	Exec: exec,

	Params: butler.NewParams(paramType),
}

func exec(r *butler.Butler, files []string) error {
	if len(files) == 0 {
		return fmt.Errorf("at least one file must be specified")
	}

	path := files[0]
	typ := r.Params.Get("type").Str()
	switch typ {
	case "":
		return errors.New("empty parse type")
	case "package":
		return parsePkg(path)
	case "unit":
		return parseUnit(path)
	default:
		return fmt.Errorf("unknown parse type: %s", typ)
	}
}

func parseUnit(path string) error {
	pool := sm.New()
	text, err := pool.Load(path)
	if err != nil {
		return err
	}

	p := parser.FromText(text)
	unit, err := p.Unit()
	if err != nil {
		return diag.Format(pool, err.(diag.Error))
	}
	_ = unit
	return nil
}

func parsePkg(path string) error {
	pool := sm.New()
	text, err := pool.Load(path)
	if err != nil {
		return err
	}

	p := parser.FromText(text)
	pkg, err := p.Package()
	if err != nil {
		return diag.Format(pool, err.(diag.Error))
	}
	for _, m := range pkg.Modules {
		links := listLinks(m.Links)
		if m.Main != nil {
			fmt.Printf("(exe) %s: %s + %v\n", m.Name.Val, m.Main.Val, links)
		} else {
			fmt.Printf("(obj) %s: %v + %v\n", m.Name.Val, listUnits(m.Units), links)
		}
	}
	return nil
}

func listUnits(units []ast.UnitEntry) []string {
	if len(units) == 0 {
		return nil
	}

	ll := make([]string, 0, len(units))
	for _, u := range units {
		ll = append(ll, u.Val)
	}
	return ll
}

func listLinks(links []ast.LinkEntry) []string {
	if len(links) == 0 {
		return nil
	}

	ll := make([]string, 0, len(links))
	for _, l := range links {
		ll = append(ll, l.Val)
	}
	return ll
}
