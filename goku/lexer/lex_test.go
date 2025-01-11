package lexer

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mebyus/ku/goku/source"
	"github.com/mebyus/ku/goku/token"
)

func TestLex(t *testing.T) {
	entries, err := os.ReadDir("testdata")
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
		path := filepath.Join("testdata", file)
		text, err := pool.Load(path)
		if err != nil {
			t.Error(err)
			return
		}

		t.Run(file, func(t *testing.T) {
			wantListFile, err := os.Open(path + ".tokens")
			if err != nil {
				t.Error(err)
				return
			}

			lx := FromText(text)
			sc := bufio.NewScanner(wantListFile)
			var line uint32 // zero-based line number (in token list file)
			for sc.Scan() {
				wantLine := strings.TrimSpace(sc.Text())
				if wantLine == "" {
					continue
				}

				gotLine := token.FormatTokenLine(pool, lx.Lex())
				if gotLine != wantLine {
					t.Errorf("(%d)\ngot   # %s #\nwant  # %s #", line+1, gotLine, wantLine)
					return
				}

				line += 1
			}
			err = sc.Err()
			if err != nil {
				t.Error(err)
				return
			}

			last := lx.Lex()
			if last.Kind != token.EOF {
				t.Errorf("(last) Lex() got = %s, want EOF", last)
				return
			}
		})
	}
}
