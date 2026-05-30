package genc

import (
	"github.com/mebyus/ku/internal/ku/enums/scok"
	"github.com/mebyus/ku/internal/ku/stg"
)

func (g *Buffer) getName(s *stg.Symbol) string {
	if s.Scope.Kind == scok.Unit {
		return g.prefix + s.Name
	}

	return s.Name
}
