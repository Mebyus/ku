package parser

import (
	"bufio"
	"bytes"
	"cmp"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/mebyus/ku/internal/fbuf"
	"github.com/mebyus/ku/internal/ku/ast"
	"github.com/mebyus/ku/internal/ku/sx"
)

func TestParse(t *testing.T) {
	const progsDir = "src"

	progs, err := listTestProgs(progsDir)
	if err != nil {
		t.Errorf("list test programs: %v", err)
		return
	}
	if len(progs) == 0 {
		t.Error("no test programs found")
		return
	}

	pool := sx.New()
	for _, name := range progs {
		path := filepath.Join(progsDir, name)
		text, err := pool.Load(path)
		if err != nil {
			t.Errorf("load program \"%s\" source text: %v", name, err)
			return
		}

		tt := ParseText(pool, text)
		err = compareTestOutput(path, tt)
		if err != nil {
			diff, ok := err.(*lineDiffError)
			if ok {
				t.Errorf("ParseText(\"%s\") line %d diff:\n[got]  %s\n[want] %s", name, diff.line, diff.got, diff.want)
			} else {
				t.Errorf("ParseText(\"%s\") error = %v", name, err)
			}
			continue
		}
	}
}

// describes first different line between two texts as an error
type lineDiffError struct {
	got  string
	want string
	line int // line number, one based value
	ndif int // -1, 0, 1 for comparing total number of lines
}

func (e *lineDiffError) Error() string {
	return ""
}

// returns non-nil error with first diff line
func compareTestOutput(path string, text *ast.Text) error {
	want, err := readFileLines(path)
	if err != nil {
		return err
	}

	var buf fbuf.Buffer
	text.Format(&buf)
	got, err := readLines(bytes.NewReader(buf.Bytes()))
	if err != nil {
		return err
	}

	return compareLines(got, want)
}

func compareLines(got, want []string) error {
	n := min(len(got), len(want))
	ndif := cmp.Compare(len(got), len(want))
	for i := range n {
		if got[i] != want[i] {
			return &lineDiffError{
				got:  got[i],
				want: want[i],
				line: i + 1,
				ndif: ndif,
			}
		}
	}

	if ndif != 0 {
		if len(got) < len(want) {
			return &lineDiffError{
				got:  "",
				want: want[n],
				line: n + 1,
				ndif: ndif,
			}
		}
		return &lineDiffError{
			got:  got[n],
			want: "",
			line: n + 1,
			ndif: ndif,
		}
	}

	return nil
}

func readFileLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return readLines(file)
}

func readLines(r io.Reader) ([]string, error) {
	s := bufio.NewScanner(r)
	var lines []string
	for s.Scan() {
		lines = append(lines, s.Text())
	}
	err := s.Err()
	if err != nil {
		return nil, err
	}
	return lines, nil
}

func listTestProgs(dir string) ([]string, error) {
	ents, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	list := make([]string, 0, len(ents))
	for _, ent := range ents {
		if ent.IsDir() || !ent.Type().IsRegular() {
			continue
		}

		name := ent.Name()
		if filepath.Ext(name) == ".ku" {
			list = append(list, name)
		}
	}
	return list, nil
}
