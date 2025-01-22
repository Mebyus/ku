package diag

import (
	"fmt"
	"io"

	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/source"
)

type UnknownOriginError struct {
	Name ast.Word
}

var _ Error = &UnknownOriginError{}

func (e *UnknownOriginError) Error() string {
	return fmt.Sprintf("unknown import origin \"%s\"", e.Name.Str)
}

func (e *UnknownOriginError) Render(w io.Writer, m source.PinMap) error {
	pos, err := m.DecodePin(e.Name.Pin)
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
