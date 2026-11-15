package main

import (
	"os"

	"github.com/mebyus/ku/internal/ku/lexer"
	"github.com/mebyus/ku/internal/ku/sx"
	"github.com/mebyus/ku/internal/ku/token"
)

func lex(path string) error {
	pool := sx.New()
	text, err := pool.Load(path)
	if err != nil {
		return err
	}

	lx := lexer.FromText(text)
	return token.Render(pool, os.Stdout, lx)
}
