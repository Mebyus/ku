package main

import (
	"fmt"
	"os"

	"github.com/mebyus/ku/goku/butler"
	"github.com/mebyus/ku/goku/cmd/ku/lex"
)

func main() {
	if len(os.Args) == 0 {
		panic("os args are empty")
	}
	args := os.Args[1:]

	err := butler.Run(root, os.Stderr, args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var root = &butler.Butler{
	Name: "ku",

	Short: "Ku is a command line tool for managing Ku source code.",
	Usage: "[command] [arguments]",

	Subs: []*butler.Butler{
		lex.Butler,
	},
}
