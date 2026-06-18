package main

import (
	"fmt"
	"os"
	"time"

	"github.com/mebyus/ku/goku/compiler/cc"
	"github.com/mebyus/ku/goku/compiler/enums/bk"
	"github.com/mebyus/ku/internal/ku/ast"
	"github.com/mebyus/ku/internal/ku/genc"
	"github.com/mebyus/ku/internal/ku/parser"
	"github.com/mebyus/ku/internal/ku/stg"
	"github.com/mebyus/ku/internal/ku/sx"
)

func build(paths []string) error {
	const debug = true

	pool := sx.New()
	texts, err := parsePaths(pool, paths)
	if err != nil {
		return err
	}

	start := time.Now()
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
	if debug {
		fmt.Printf("parse: %s\n", time.Since(start))
	}

	tp := stg.NewPool(pool)
	t := tp.Get()

	start = time.Now()
	unit := t.Translate(parsed)
	if debug {
		fmt.Printf("stg:   %s\n", time.Since(start))
	}

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

	const progName = "out.c"
	start = time.Now()
	err = genProg(progName, &prog)
	if err != nil {
		fmt.Fprintf(os.Stderr, "genc: %s\n", err)
		os.Exit(1)
	}
	if debug {
		fmt.Printf("genc:  %s\n", time.Since(start))
	}

	start = time.Now()
	err = cc.CompileObj("out.o", progName, bk.Debug)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cc: %s\n", err)
		os.Exit(1)
	}
	if debug {
		fmt.Printf("cc:    %s\n", time.Since(start))
	}

	return nil
}

func genProg(out string, prog *stg.Program) error {
	file, err := os.Create(out)
	if err != nil {
		return err
	}
	defer file.Close()

	return genc.Gen(file, prog)
}
