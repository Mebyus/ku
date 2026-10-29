package genc_test

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/mebyus/ku/internal/ku/builder"
)

// TestRun is an end-to-end test of Ku compiler and
// its C backend (genc).
//
// Test imitates real-world scenario when program is compiled and
// saved into native executable format and then executed. In more detail
// test consists of the following steps:
//
//  1. Compile program from a Ku source file into C code.
//  2. Compile resulting C code via C compiler into executable.
//  3. Execute with given arguments.
//  4. Compare output (from stdout) with expected test result.
//
// Each program should be executed multiple times with different
// arguments (expected results change based on given arguments).
//
// This test is positive. All programs used here must compile
// and execute without errors.
//
// Test program sources are stored in "./src" directory.
//
// Function called "run" in Ku source acts as entrypoint in
// this test.
func TestRun(t *testing.T) {
	bundles := []Bundle{
		{
			// empty run function
			src: "0000_empty_run.ku",
			tests: []TestCase{
				{
					args: nil,
					want: nil,
				},
			},
		},
	}

	// validate bundles and test cases
	n := 0 // number of validation errors
	set := make(map[string]struct{}, len(bundles))
	for i, b := range bundles {
		name := b.src
		if name == "" {
			t.Errorf("bundle (i=%d) has no source file", i)
			n += 1
		} else {
			_, ok := set[name]
			if ok {
				t.Errorf("duplicate source file \"%s\" in bundles", name)
				n += 1
			}
			set[name] = struct{}{}
		}

		if len(b.tests) == 0 {
			t.Errorf("no test cases in \"%s\" bundle", name)
			n += 1
		}
	}
	if n != 0 {
		return
	}

	dir := t.TempDir()
	root, err := getRoot()
	if err != nil {
		t.Errorf("get lang root directory: %v", err)
		return
	}

	t.Logf("temp dir: %s", dir)
	t.Logf("lang dir: %s", root)

	var out bytes.Buffer
	for _, b := range bundles {
		name := b.src
		path := filepath.Join("src", name)

		exe, err := builder.CompileFileForTestGenC(dir, path)
		if err != nil {
			t.Errorf("compile program \"%s\": %v", name, err)
			continue
		}

		for _, c := range b.tests {
			got, err := runExe(t.Context(), exe, &out, c.args)
			if err != nil {
				t.Errorf("run test executable \"%s\": %v", name, err)
				continue
			}

			err = compareOutputLines(got, c.want)
			if err != nil {
				t.Errorf("compare output \"%s\" for args = %v: %v", name, c.args, err)
				for i, line := range got {
					t.Logf("%d: %s", i+1, line)
				}
			}
		}
	}
}

// runs executable file specified by path and returns its
// output (from stdout) in form of separated text lines
// (without \n endings)
func runExe(ctx context.Context, path string, buf *bytes.Buffer, args []string) ([]string, error) {
	buf.Reset()

	const timeout = 250 * time.Millisecond

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, path, args...)
	cmd.Stdout = buf

	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(buf)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	err = scanner.Err()
	if err != nil {
		return nil, err
	}

	return lines, nil
}

// returns path to repository root directory
//
// workaround for Go's strange policy for setting current directory
// inside "go test" command
func getRoot() (string, error) {
	const n = 10 // max iterations

	dir := ""
	for range n {
		dir += "../"
		path := filepath.Join(dir, "go.mod")
		_, err := os.Stat(path)
		if err == nil {
			return dir, nil
		}
		if !os.IsNotExist(err) {
			return "", err
		}
	}

	return "", os.ErrNotExist
}

// returns non-nil error with description if provided
// slices are not equal
func compareOutputLines(a, b []string) error {
	if len(a) != len(b) {
		return fmt.Errorf("line number mismatch (got %d, want %d)", len(a), len(b))
	}

	n := len(a)
	for i := range n {
		if a[i] != b[i] {
			return fmt.Errorf("different line %d (want \"%s\")", i+1, b[i])
		}
	}

	return nil
}

type TestCase struct {
	// program arguments
	args []string

	// expected output (from stdout) after execution
	//
	// contains expected text lines, without \n endings
	want []string
}

// Bundle is a collection of test cases for a single program.
type Bundle struct {
	tests []TestCase

	// filename of source program
	src string
}
