package diag

import (
	"io"

	"github.com/mebyus/ku/goku/compiler/sm"
)

// SimpleMessageError stores a simple message with position information.
// This error is rendered as single text line in output.
type SimpleMessageError struct {
	Text string
	Pin  sm.Pin
}

var _ Error = &SimpleMessageError{}

func (e *SimpleMessageError) Error() string {
	return e.Text
}

func (e *SimpleMessageError) Render(w io.Writer, m sm.PinMap) error {
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

func (e *SimpleMessageError) SetFallbackSpan(span sm.Span) {
	if e.Pin == 0 {
		e.Pin = span.Pin
	}
}
