package builder

import "github.com/mebyus/ku/goku/compiler/srcmap"

type Unit struct {
	// Order of elements directly corresponds to file include order
	// in unit build file.
	Texts []*srcmap.Text

	// List of test functions found in unit source files.
	Tests []string

	// Does not include unit or main directory prefix.
	Path string

	Main bool
}

type Module struct {
	Units []*Unit

	Name string

	Main *Unit
}
