package genc

import (
	"strings"

	"github.com/mebyus/ku/internal/ku/enums/scok"
	"github.com/mebyus/ku/internal/ku/enums/symk"
	"github.com/mebyus/ku/internal/ku/stg"
)

func (g *genbuf) getName(s *stg.Symbol) string {
	if s.Link != "" {
		return s.Link
	}
	if s.Kind == symk.Method {
		return g.prefix + "m_" + strings.Replace(s.Name, ".", "_", 1)
	}

	if s.Scope.Kind == scok.Unit {
		if s.Kind == symk.Fun {
			return g.prefix + "f_" + s.Name
		}
		if s.Kind == symk.Type {
			return g.prefix + "t_" + s.Name
		}
		return g.prefix + s.Name
	}

	return s.Name
}
