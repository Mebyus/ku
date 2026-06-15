package main

import (
	"fmt"
	"os"

	"github.com/mebyus/ku/internal/ku/ast"
	"github.com/mebyus/ku/internal/ku/genc"
	"github.com/mebyus/ku/internal/ku/parser"
	"github.com/mebyus/ku/internal/ku/stg"
	"github.com/mebyus/ku/internal/ku/sx"
)

func build(paths []string) error {
	pool := sx.New()
	texts, err := parsePaths(pool, paths)
	if err != nil {
		return err
	}

	n := 0 // total number of errors
	var parsed []*ast.Text
	for _, x := range texts {
		t := parser.ParseText(x)
		for _, e := range t.Errors {
			pos := pool.DecodePin(e.Pin)
			fmt.Fprintf(os.Stderr, "%s: %s\n", pos, e.Short)
			n += 1
		}
		parsed = append(parsed, t)
	}

	tp := stg.NewPool(pool)
	t := tp.Get()
	unit := t.Translate(parsed)
	var prog stg.Program

	for _, e := range unit.Errors {
		pos := pool.DecodePin(e.Pin)
		fmt.Fprintf(os.Stderr, "%s: %s\n", pos, e.Short)
		n += 1
	}

	if n != 0 {
		os.Exit(1)
	}

	prog.Common = tp.Common
	prog.Units = []*stg.Unit{unit}
	return genc.Gen(os.Stdout, &prog)
}
