package diag

import (
	"fmt"
	"io"

	"github.com/mebyus/ku/goku/compiler/srcmap"
	"github.com/mebyus/ku/goku/compiler/token"
)

// ParseError describes an error which occured during source code parsing phase.
type ParseError struct {
	Text string
	Pin  srcmap.Pin
}

type UnexpectedTokenError struct {
	Token token.Token
}

var _ Error = &UnexpectedTokenError{}

func (e *UnexpectedTokenError) Error() string {
	return fmt.Sprintf("unexpected token \"%s\"", e.Token.String())
}

func (e *UnexpectedTokenError) Render(w io.Writer, m srcmap.PinMap) error {
	pos, err := m.DecodePin(e.Token.Pin)
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
