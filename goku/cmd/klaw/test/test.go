package test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/mebyus/ku/goku/butler"
	"github.com/mebyus/ku/goku/compiler/cc"
	"github.com/mebyus/ku/goku/compiler/enums/bk"
	"github.com/mebyus/ku/goku/klaw/builder"
)

var Butler = &butler.Butler{
	Name: "test",

	Short: "Test Ku unit (via C codegen)",
	Usage: "[options] [path]",

	Params: butler.NewParams(
		butler.Param{
			Name:    "build-kind",
			Alias:   "k",
			Desc:    "Specifies build kind",
			Default: bk.Debug.String(),
			Kind:    butler.String,
		},
	),

	Exec: run,
}

func run(r *butler.Butler, units []string) error {
	if len(units) == 0 {
		return fmt.Errorf("at least one unit must be specified")
	}

	path := units[0]
	return test(path, r.Params.Get("build-kind").Str())
}

func test(path string, kind string) error {
	k, err := bk.Parse(kind)
	if err != nil {
		return err
	}

	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	if !info.IsDir() {
		return fmt.Errorf("\"%s\" is not a directory", path)
	}

	codegenOutPath := getCodegenOutPath(path)
	err = mkdir(filepath.Dir(codegenOutPath))
	if err != nil {
		return err
	}

	err = genFromUnit(codegenOutPath, path)
	if err != nil {
		return err
	}

	testExePath := getTestExecutablePath(path)
	err = mkdir(filepath.Dir(testExePath))
	if err != nil {
		return err
	}

	err = cc.CompileExe(testExePath, codegenOutPath, k)
	if err != nil {
		return err
	}
	return runTestExecutable(testExePath)
}

func runTestExecutable(path string) error {
	cmd := exec.Command(path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func genFromUnit(out, path string) error {
	genOut, err := os.Create(out)
	if err != nil {
		return err
	}
	defer genOut.Close()

	return builder.GenUnitWithTests(genOut, path)
}

func getCodegenOutPath(path string) string {
	return filepath.Join(".kubout/genc", filepath.Base(path)+".test.kubgen.c")
}

func getTestExecutablePath(path string) string {
	return filepath.Join(".kubout/test", filepath.Base(path))
}

func mkdir(path string) error {
	if path == "" || path == "." {
		return nil
	}
	return os.MkdirAll(path, 0o755)
}
