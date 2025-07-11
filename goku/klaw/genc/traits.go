package genc

import "github.com/mebyus/ku/goku/compiler/ast"

func getPropValue(traits ast.Traits, name string) (string, bool) {
	if traits.Props == nil {
		return "", false
	}

	props := *traits.Props
	for _, p := range props {
		if p.Name == name {
			return p.Exp.(ast.String).Val, true
		}
	}
	return "", false
}
