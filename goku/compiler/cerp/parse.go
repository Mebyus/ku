package cerp

import "github.com/mebyus/ku/goku/compiler/source"

// ParseError describes an error which occured during source code parsing phase.
type ParseError struct {
	Pin source.Pin
}
