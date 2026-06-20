package main

import (
	"fmt"

	"github.com/mebyus/ku/internal/ku/builder"
)

func walk(path string) error {
	list := builder.Walk(path)
	for _, p := range list {
		fmt.Println(p)
	}
	return nil
}
