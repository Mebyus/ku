package main

import (
	"os"

	"github.com/mebyus/ku/internal/ku/builder"
)

func clean() error {
	dir, err := builder.GetOutDir()
	if err != nil {
		return err
	}
	return os.RemoveAll(dir)
}
