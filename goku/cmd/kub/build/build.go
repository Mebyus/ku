package build

import (
	"fmt"

	"github.com/mebyus/ku/goku/butler"
	"github.com/mebyus/ku/goku/compiler/enums/bk"
	"github.com/mebyus/ku/goku/kub/builder"
)

var Butler = &butler.Butler{
	Name: "build",

	Short: "Build (via C codegen) specified module from Kub package",
	Usage: "[options] [path]",

	Params: butler.NewParams(
		butler.Param{
			Name:    "out",
			Alias:   "o",
			Desc:    "Path to output file",
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
	),

	Exec: exec,
}

func exec(r *butler.Butler, modules []string) error {
	if len(modules) == 0 {
		return fmt.Errorf("at least one module must be specified")
	}

	name := modules[0]
	return build(name, r.Params.Get("out").Str(), r.Params.Get("build-kind").Str())
}

func build(name string, out string, kind string) error {
	k, err := bk.Parse(kind)
	if err != nil {
		return err
	}

	return builder.BuildModule(&builder.BuildModuleConfig{
		Config: builder.Config{
			ModuleName:      name,
			BuildKind:       k,
			PackageFilePath: "pkg.kub",
		},
		OutputPath: out,
	})
}
