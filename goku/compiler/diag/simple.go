package diag

import (
	"io"

	"github.com/mebyus/ku/goku/compiler/srcmap"
)

// SimpleMessageError stores a simple message with position information.
// This error is rendered as single text line in output.
type SimpleMessageError struct {
	Text string
	Pin  srcmap.Pin
}

var _ Error = &SimpleMessageError{}

func (e *SimpleMessageError) Error() string {
	return e.Text
}

func (e *SimpleMessageError) Render(w io.Writer, m srcmap.PinMap) error {
	pos, err := m.DecodePin(e.Pin)
	if err != nil {
		return err
	}
	_, err = io.WriteString(w, pos.String())
	if err != nil {
		return err
	}
	_, err = io.WriteString(w, " ")
	if err != nil {
		return err
	}
	_, err = io.WriteString(w, e.Error())
	return err
}
