package ast

import "github.com/mebyus/ku/internal/ku/sx"

// Import represents import of a single unit under specified symbol name.
type Import struct {
	// Unit path of imported unit.
	Path sx.Path

	// Name under which unit is being imported.
	Name string

	// Name pin.
	Pin sx.Pin

	// Pin of import string.
	ImpPin sx.Pin
}
