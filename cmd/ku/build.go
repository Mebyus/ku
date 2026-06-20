package main

import (
	"github.com/mebyus/ku/internal/ku/builder"
)

func build(path string) error {
	r := builder.Build(&builder.Config{Unit: path})
	return r.Error
}
