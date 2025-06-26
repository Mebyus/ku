package diag

import (
	"io"

	"github.com/mebyus/ku/goku/compiler/srcmap"
)

type PinlessError struct {
	Text string
}

var _ Error = &PinlessError{}

func (e *PinlessError) Error() string {
	return e.Text
}

func (e *PinlessError) Render(w io.Writer, m srcmap.PinMap) error {
	_, err := io.WriteString(w, e.Text)
	return err
}
