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
		path := os.Args[2]
		err = parse(path)
	case "build":
		path := os.Args[2]
		err = build(path)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func parseFromPath(path string) (*sx.Pool, []*sx.Text, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, nil, err
	}

	pool := sx.New()
	var texts []*sx.Text
	if info.IsDir() {
		list, err := pool.LoadDir(&sx.LoadParams{
			Dir:              path,
			IncludeTestFiles: true,
		})
		if err != nil {
			return nil, nil, err
		}
		texts = list
	} else {
		text, err := pool.Load(path)
		if err != nil {
			return nil, nil, err
		}
		texts = append(texts, text)
	}

	return pool, texts, nil
}

func parse(path string) error {
	pool, texts, err := parseFromPath(path)
	if err != nil {
		return err
	}

	n := 0 // total number of errors
	for _, x := range texts {
		t := parser.ParseText(x)
		for _, e := range t.Errors {
			pos := pool.DecodePin(e.Pin)
			fmt.Fprintf(os.Stderr, "%s: %s\n", pos, e.Short)
			n += 1
		}
	}

	if n != 0 {
		os.Exit(1)
	}

	return nil
}
