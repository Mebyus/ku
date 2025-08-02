package eval

import (
	"github.com/mebyus/ku/goku/compiler/enums/bk"
	"github.com/mebyus/ku/goku/compiler/enums/bm"
)

type Env struct {
	BuildKind bk.Kind
	BuildMode bm.Mode

	m map[string]Value
}

func NewEnv() *Env {
	return &Env{m: make(map[string]Value)}
}
