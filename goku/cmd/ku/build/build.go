package build

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/mebyus/ku/goku/butler"
	"github.com/mebyus/ku/goku/compiler/builder"
	"github.com/mebyus/ku/goku/compiler/enums/bk"
)

var Butler = &butler.Butler{
	Name: "build",

	Short: "Build specified Ku unit into one object file",
	Usage: "[options] [unit]",

	Params: butler.NewParams(
		butler.Param{
			Name:    "out",
			Alias:   "o",
			Desc:    "Path to output object file",
			Default: "",
			Kind:    butler.String,
		},
		butler.Param{
			Name:    "build-kind",
			Alias:   "k",
			Desc:    "Specifies build kind",
			Default: bk.Debug.String(),
			Kind:    butler.String,
		},
		butler.Param{
			Name:    "cc-include-dirs",
			Desc:    "List of additional include directories for C compiler",
			Default: []string{},
			Kind:    butler.List,
		},
	),

	Exec: exec,
}

func exec(r *butler.Butler, list []string) error {
	if len(list) == 0 {
		return errors.New("unit must be specified")
	}
	if len(list) != 1 {
		return errors.New("only one unit may be specified")
	}

	unit := strings.TrimSpace(list[0])
	if unit == "" {
		return errors.New("empty unit path")
	}

	// We accept unit path in two forms:
	// 1. with "./src/" prefix which means we a dealing with filesystem path to directory
	// 2. relative path without leading ".", then we treat it as unit path relative to source directory
	unit = strings.TrimPrefix(unit, "./src/")

	kind, err := bk.Parse(r.Params.Get("build-kind").Str())
	if err != nil {
		return err
	}

	out := r.Params.Get("out").Str()
	if out == "." || out == ".." {
		return errors.New("invalid output path")
	}
	if out != "" {
		out = filepath.Clean(out)
	}

	return builder.Build(&builder.Config{
		Unit:      unit,
		OutPath:   out,
		BuildKind: kind,
	})
}
