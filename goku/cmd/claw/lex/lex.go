package lex

import (
	"fmt"
	"os"

	"github.com/mebyus/ku/goku/butler"
	"github.com/mebyus/ku/goku/claw/lexer"
	"github.com/mebyus/ku/goku/compiler/srcmap"
)

var Butler = &butler.Butler{
	Name: "lex",

	Short: "list token stream produced by a given source file",
	Usage: "[options] [file]",

	Exec: exec,
}

func exec(r *butler.Butler, files []string) error {
	if len(files) == 0 {
		return fmt.Errorf("at least one file must be specified")
	}
	return lex(files[0])
}

func lex(path string) error {
	pool := srcmap.New()
	text, err := pool.Load(path)
	if err != nil {
		return err
	}

	lx := lexer.FromText(text)
	return lexer.Render(os.Stdout, lx, pool)
}
