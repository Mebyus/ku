package sx

import (
	"fmt"
	"io"
)

// Error describes compiler error with source position.
type Error struct {
	// Short error message that describes what is wrong.
	Short string

	// Possibly long error description and additional info
	// with meaning varying based on type of error.
	//
	// May provide user help, suggestions, advice on how to fix.
	Note string

	// What place (in program source code) error should be attributed.
	//
	// Equals 0 if such place cannot be specified.
	Pin Pin
}

func FormatError(pool *Pool, out io.Writer, e *Error) {
	pos := pool.FormatPin(e.Pin)
	if pos == "" {
		fmt.Fprintf(out, "%s\n", e.Short)
	} else {
		fmt.Fprintf(out, "%s: %s\n", pos, e.Short)
	}

	if e.Note != "" {
		fmt.Fprintf(out, "======\n%s\n======\n", e.Note)
	}
}
