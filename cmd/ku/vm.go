package main

import (
	"strings"

	"github.com/mebyus/ku/internal/ku/builder"
)

func vmc(path string) error {
	out := strings.TrimSuffix(path, ".ku") + ".kvx"
	return builder.CompileFileForVM32B(path, out)
}
