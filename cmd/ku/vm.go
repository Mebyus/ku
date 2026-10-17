package main

import "github.com/mebyus/ku/internal/ku/builder"

func vmc(path string) error {
	return builder.CompileFileForVM32B(path)
}
