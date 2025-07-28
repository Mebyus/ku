package test

import (
	"fmt"

	"github.com/mebyus/ku/goku/butler"
	"github.com/mebyus/ku/goku/compiler/enums/bk"
	"github.com/mebyus/ku/goku/kub/builder"
)

var Butler = &butler.Butler{
	Name: "test",

	Short: "Test Kub module (via C codegen)",
	Usage: "[options] [path]",

	Params: butler.NewParams(
		butler.Param{
			Name:    "build-kind",
			Alias:   "k",
			Desc:    "Specifies build kind",
			Default: bk.Debug.String(),
			Kind:    butler.String,
		},
		butler.Param{
			Name:    "name",
			Desc:    "Specifies test name",
			Default: "",
			Kind:    butler.String,
		},
	),

	Exec: run,
}

func run(r *butler.Butler, units []string) error {
	if len(units) == 0 {
		return fmt.Errorf("at least one module must be specified")
	}

	path := units[0]
	return test(path, r.Params.Get("build-kind").Str(), r.Params.Get("name").Str())
}

func test(name string, kind string, test string) error {
	k, err := bk.Parse(kind)
	if err != nil {
		return err
	}

	return builder.TestModule(&builder.TestModuleConfig{
		Config: builder.Config{
			ModuleName:      name,
			BuildKind:       k,
			PackageFilePath: "pkg.kub",
		},
		TestName: test,
	})
}
