package main

import (
	"fmt"
	"os"

	"github.com/mebyus/ku/goku/butler"
	"github.com/mebyus/ku/goku/cmd/kub/compile"
	"github.com/mebyus/ku/goku/cmd/kub/eval"
	"github.com/mebyus/ku/goku/cmd/kub/genc"
	"github.com/mebyus/ku/goku/cmd/kub/lex"
	"github.com/mebyus/ku/goku/cmd/kub/parse"
	"github.com/mebyus/ku/goku/cmd/kub/run"
	"github.com/mebyus/ku/goku/cmd/kub/test"
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
	Name: "kub",

	Short: "Kub is a command line tool for managing restricted set of Ku source code.",
	Usage: "[command] [arguments]",

	Subs: []*butler.Butler{
		lex.Butler,
		parse.Butler,
		eval.Butler,
		genc.Butler,
		compile.Butler,
		test.Butler,
		run.Butler,
	},
}
