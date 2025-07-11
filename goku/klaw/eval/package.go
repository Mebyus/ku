package eval

import (
	"fmt"
	"strings"

	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/klaw/ast"
)

type Package struct {
	Modules []Module

	MainDir   string
	SourceDir string
}

type Module struct {
	Links []string

	// List of root units.
	Units []string

	// Empty for object module.
	Main string

	Name string
}

func EvalPackage(pkg *ast.Package) (*Package, diag.Error) {
	var p Package
	err := evalSets(&p, pkg.Sets)
	if err != nil {
		return nil, err
	}
	err = evalModules(&p, pkg.Modules)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func evalModules(pkg *Package, modules []ast.Module) diag.Error {
	for _, m := range modules {
		err := evalModule(pkg, m)
		if err != nil {
			return err
		}
	}
	return nil
}

func evalModule(pkg *Package, module ast.Module) diag.Error {
	var mainName string
	if module.Main != nil {
		mainName = module.Main.Val
	}

	m := Module{
		Name:  module.Name.Val,
		Main:  mainName,
		Units: listUnits(module.Units),
		Links: listLinks(module.Links),
	}
	if m.Main == "" && len(m.Units) == 0 {
		return &diag.PinlessError{Text: fmt.Sprintf("no units in \"%s\" module", m.Name)}
	}
	pkg.Modules = append(pkg.Modules, m)
	return nil
}

func evalSets(pkg *Package, sets []ast.Set) diag.Error {
	for _, s := range sets {
		err := evalSet(pkg, s)
		if err != nil {
			return err
		}
	}
	return nil
}

func evalSet(pkg *Package, set ast.Set) diag.Error {
	name := joinNameParts(set.Name.Parts)
	switch name {
	case "source.dir":
	case "main.dir":
	}
	return nil
}

func listLinks(entries []ast.LinkEntry) []string {
	if len(entries) == 0 {
		return nil
	}
	list := make([]string, 0, len(entries))
	for _, l := range entries {
		list = append(list, l.Val)
	}
	return list

}

func listUnits(entries []ast.UnitEntry) []string {
	if len(entries) == 0 {
		return nil
	}
	list := make([]string, 0, len(entries))
	for _, u := range entries {
		list = append(list, u.Val)
	}
	return list
}

func joinNameParts(parts []ast.Word) string {
	if len(parts) == 0 {
		panic("no parts")
	}

	list := make([]string, 0, len(parts))
	for _, p := range parts {
		list = append(list, p.Str)
	}
	return strings.Join(list, ".")
}
