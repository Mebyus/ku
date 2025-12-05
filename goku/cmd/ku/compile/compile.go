package compile

import (
	"errors"
	"path/filepath"

	"github.com/mebyus/ku/goku/butler"
	"github.com/mebyus/ku/goku/compiler/enums/bk"
	"github.com/mebyus/ku/goku/compiler/srcmap"
	"github.com/mebyus/ku/goku/kub/builder"
)

var Butler = &butler.Butler{
	Name: "compile",

	Short: "Compile one or more Ku or C source files into one object file",
	Usage: "[options] [files]",

	Params: butler.NewParams(
		butler.Param{
			Name:  "out",
			Alias: "o",
			Desc:  "Path to output object file",
			Kind:  butler.String,
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

func exec(r *butler.Butler, files []string) error {
	if len(files) == 0 {
		return errors.New("at least one source file must be specified")
	}
	kind, err := bk.Parse(r.Params.Get("build-kind").Str())
	if err != nil {
		return err
	}
	out := filepath.Clean(r.Params.Get("out").Str())
	if out == "" || out == "." || out == ".." {
		return errors.New("empty or invalid output path")
	}

	return compile(&builder.CompileConfig{
		IncludeDirs: r.Params.Get("cc-include-dirs").List(),
		OutputPath:  out,
		BuildKind:   kind,
	}, files)
}

func compile(c *builder.CompileConfig, files []string) error {
	pool := srcmap.New()

	texts := make([]*srcmap.Text, 0, len(files))
	for _, path := range files {
		text, err := pool.Load(path)
		if err != nil {
			return err
		}
		texts = append(texts, text)
	}

	return builder.CompileTexts(c, texts)
}
