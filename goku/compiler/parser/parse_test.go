package parser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/lexer"
	"github.com/mebyus/ku/goku/compiler/source"
)

// Idea of this test: compare token streams in the following round trip.
//
//	Load("example.ku") => source.Text (1)
//	Lex(Text) => Stream (1)
//	Parse(Stream) => ast.Text
//	Print(Text) => source.Text (2)
//	Lex(Text) => Stream (2)
//
// Streams 1 and 2 should be identical to each other (if parsing was successful).
// The only exception could be source texts which include dangling comma in some
// cases (struct fields, function parameters and args, etc.). We will avoid such
// texts in test data.
func TestParse(t *testing.T) {
	const dir = "../ast/testdata"
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Error(err)
		return
	}

	var files []string
	for _, entry := range entries {
		name := entry.Name()
		if strings.HasSuffix(name, ".ku") {
			files = append(files, name)
		}
	}

	pool := source.New()
	for _, file := range files {
		path := filepath.Join(dir, file)
		src, err := pool.Load(path)
		if err != nil {
			t.Error(err)
			return
		}

		t.Run(file, func(t *testing.T) {
			tokens := lexer.ListTokens(lexer.FromText(src))
			lx1 := lexer.FromTokens(tokens)
			text, err := ParseStream(lx1)
			if err != nil {
				t.Error(diag.Stringify(pool, err))
				return
			}

			var printer ast.Printer
			printer.Text(text)

			lx1.Reset()
			lx2 := lexer.FromBytes(printer.Bytes())
			c := lexer.Compare(lx2, lx1)
			if c != nil {
				t.Error(c)
				return
			}
		})
	}
}
