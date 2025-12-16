package diag

import (
	"fmt"
	"io"

	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/srcmap"
)

type UnknownOriginError struct {
	Name ast.Word
}

var _ Error = &UnknownOriginError{}

func (e *UnknownOriginError) Error() string {
	return fmt.Sprintf("unknown import origin \"%s\"", e.Name.Str)
}

func (e *UnknownOriginError) Render(w io.Writer, m srcmap.PinMap) error {
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

type ImportCycleError struct {
	Sites []srcmap.ImportSite
}

var _ Error = &ImportCycleError{}

func (e *ImportCycleError) Error() string {
	return "import cycle detected"
}

func (e *ImportCycleError) Render(w io.Writer, m srcmap.PinMap) error {
	_, err := io.WriteString(w, "import cycle detected:\n")
	if err != nil {
		return err
	}
	for _, s := range e.Sites {
		err = renderImport(w, m, s)
		if err != nil {
			return err
		}
	}
	return nil
}

func renderImport(w io.Writer, m srcmap.PinMap, site srcmap.ImportSite) error {
	pos, err := m.DecodePin(site.Pin)
	if err != nil {
		return err
	}
	_, err = io.WriteString(w, "    ")
	if err != nil {
		return err
	}
	_, err = io.WriteString(w, pos.String())
	if err != nil {
		return err
	}
	_, err = io.WriteString(w, ": ")
	if err != nil {
		return err
	}
	_, err = io.WriteString(w, site.Name)
	if err != nil {
		return err
	}
	_, err = io.WriteString(w, " => \"")
	if err != nil {
		return err
	}
	_, err = io.WriteString(w, site.Path.String())
	if err != nil {
		return err
	}
	_, err = io.WriteString(w, "\"\n")
	if err != nil {
		return err
	}
	return nil
}
