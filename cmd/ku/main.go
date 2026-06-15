package main

import (
	"fmt"
	"os"

	"github.com/mebyus/ku/internal/ku/parser"
	"github.com/mebyus/ku/internal/ku/sx"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "target file of directory not specified\n")
		os.Exit(1)
	}

	var err error
	switch os.Args[1] {
	case "parse":
		paths := os.Args[2:]
		err = parse(paths)
	case "build":
		paths := os.Args[2:]
		err = build(paths)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func parsePath(pool *sx.Pool, path string) ([]*sx.Text, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	var texts []*sx.Text
	if info.IsDir() {
		list, err := pool.LoadDir(&sx.LoadParams{
			Dir:              path,
			IncludeTestFiles: true,
		})
		if err != nil {
			return nil, err
		}
		texts = list
	} else {
		text, err := pool.Load(path)
		if err != nil {
			return nil, err
		}
		texts = append(texts, text)
	}

	return texts, nil
}

func parsePaths(pool *sx.Pool, paths []string) ([]*sx.Text, error) {
	var texts []*sx.Text
	for _, p := range paths {
		ts, err := parsePath(pool, p)
		if err != nil {
			return nil, err
		}
		texts = append(texts, ts...)
	}
	return texts, nil
}

func parse(paths []string) error {
	pool := sx.New()
	texts, err := parsePaths(pool, paths)
	if err != nil {
		return err
	}

	var funs []string
	var stubs []string

	n := 0 // total number of errors
	for _, x := range texts {
		t := parser.ParseText(x)
		for _, e := range t.Errors {
			pos := pool.DecodePin(e.Pin)
			fmt.Fprintf(os.Stderr, "%s: %s\n", pos, e.Short)
			n += 1
		}

		for _, f := range t.Funs {
			funs = append(funs, f.Name)
		}

		for _, s := range t.Stubs {
			stubs = append(stubs, s.Name)
		}
	}

	if n != 0 {
		os.Exit(1)
	}

	fmt.Printf("funs:  %v\n", funs)
	fmt.Printf("stubs: %v\n", stubs)

	return nil
}
