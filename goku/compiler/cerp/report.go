package cerp

import (
	"io"

	"github.com/mebyus/ku/goku/compiler/source"
)

// Report combines several compilation errors into one report.
type Report struct {
	Errors []Error
}

// Error augments standard Go error with source position information.
//
// This interface is a container for any error related to compilation
// due to source code issues.
//
// The majority of such errors fall into syntax or semantic category.
// Syntax errors are detected during parsing. Semantic errors are detected
// in typechecking phase.
type Error interface {
	error

	Render(w io.Writer, m source.PinMap) error
}

func Render(w io.Writer, m source.PinMap, r Report) error {
	for _, e := range r.Errors {
		err := e.Render(w, m)
		if err != nil {
			return err
		}
	}
	return nil
}
