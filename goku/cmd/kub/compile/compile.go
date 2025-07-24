package compile

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mebyus/ku/goku/butler"
	"github.com/mebyus/ku/goku/compiler/cc"
	"github.com/mebyus/ku/goku/compiler/enums/bk"
	"github.com/mebyus/ku/goku/kub/builder"
)

var Butler = &butler.Butler{
	Name: "compile",

	Short: "Compile Ku unit (via C codegen)",
	Usage: "[options] [path]",

	Params: butler.NewParams(
		butler.Param{
			Name:  "out",
			Alias: "o",
			Desc:  "Path to output file (object)",
			Kind:  butler.String,
		},
		butler.Param{
			Name:    "build-kind",
			Alias:   "k",
			Desc:    "Specifies build kind",
			Default: bk.Debug.String(),
			Kind:    butler.String,
		},
	),

	Exec: exec,
}

func exec(r *butler.Butler, units []string) error {
	if len(units) == 0 {
		return fmt.Errorf("at least one unit must be specified")
	}

	path := units[0]
	return compile(r.Params.Get("out").Str(), path, r.Params.Get("build-kind").Str())
}

func compile(out, path string, kind string) error {
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

	err = mkdir(filepath.Dir(out))
	if err != nil {
		return err
	}

	return cc.CompileObj(out, codegenOutPath, k)
}

func genFromUnit(out, path string) error {
	genOut, err := os.Create(out)
	if err != nil {
		return err
	}
	defer genOut.Close()

	return builder.GenUnit(genOut, path)
}

func getCodegenOutPath(path string) string {
	return filepath.Join(".kubout/genc", filepath.Base(path)+".kubgen.c")
}

func mkdir(path string) error {
	if path == "" || path == "." {
		return nil
	}
	return os.MkdirAll(path, 0o755)
}
