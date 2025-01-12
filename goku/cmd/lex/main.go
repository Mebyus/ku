package main

import (
	"fmt"
	"os"

	"github.com/mebyus/ku/goku/compiler/lexer"
	"github.com/mebyus/ku/goku/compiler/source"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "file not specified")
		os.Exit(1)
	}

	pool := source.New()
	text, err := pool.Load(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	lx := lexer.FromText(text)
	err = lexer.ListTokens(os.Stdout, lx, pool)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
