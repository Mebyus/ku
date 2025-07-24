package run

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/mebyus/ku/goku/butler"
	"github.com/mebyus/ku/goku/compiler/cc"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/enums/bk"
	"github.com/mebyus/ku/goku/compiler/srcmap"
	"github.com/mebyus/ku/goku/kub/builder"
	"github.com/mebyus/ku/goku/kub/eval"
	"github.com/mebyus/ku/goku/kub/parser"
)

var Butler = &butler.Butler{
	Name: "run",

	Short: "Build (via C codegen) and run specified module from Kub package",
	Usage: "[options] [name]",

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

func run(r *butler.Butler, names []string) error {
	if len(names) == 0 {
		return fmt.Errorf("at least one module must be specified")
	}

	name := names[0]
	return runModule(name, r.Params.Get("build-kind").Str())
}

func runModule(name string, kind string) error {
	k, err := bk.Parse(kind)
	if err != nil {
		return err
	}

	pool := srcmap.New()
	text, err := pool.Load("pkg.kub")
	if err != nil {
		return err
	}

	p := parser.FromText(text)
	astPkg, err := p.Package()
	if err != nil {
		return diag.Format(pool, err.(diag.Error))
	}

	pkg, err := eval.EvalPackage(astPkg)
	if err != nil {
		return diag.Format(pool, err.(diag.Error))
	}

	mod := getPackageModule(pkg, name)
	if mod == nil {
		return fmt.Errorf("package has no module \"%s\"", name)
	}
	if mod.Main == "" {
		return fmt.Errorf("module \"%s\" does not specify main unit", name)
	}

	codegenOutPath := getCodegenOutPath(name)
	err = mkdir(filepath.Dir(codegenOutPath))
	if err != nil {
		return err
	}

	err = genFromMain(codegenOutPath, &builder.GenProgramConfig{
		Main:      name,
		MainDir:   pkg.MainDir,
		SourceDir: pkg.SourceDir,
		RootDir:   pkg.RootDir,
		Pool:      pool,
		BuildKind: k,
	})
	if err != nil {
		return err
	}

	exePath := getExecutablePath(name)
	err = mkdir(filepath.Dir(exePath))
	if err != nil {
		return err
	}

	err = cc.CompileExe(exePath, codegenOutPath, k)
	if err != nil {
		return err
	}
	return runExecutable(exePath)
}

func genFromMain(out string, c *builder.GenProgramConfig) error {
	genOut, err := os.Create(out)
	if err != nil {
		return err
	}
	defer genOut.Close()

	return builder.GenFromMain(genOut, c)
}

func getPackageModule(pkg *eval.Package, name string) *eval.Module {
	for _, m := range pkg.Modules {
		if m.Name == name {
			return &m
		}
	}
	return nil
}

func runExecutable(path string) error {
	cmd := exec.Command(path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func getCodegenOutPath(name string) string {
	return filepath.Join(".kubout/genc", name+".kubgen.c")
}

func getExecutablePath(name string) string {
	return filepath.Join(".kubout/exe", name)
}

func mkdir(path string) error {
	if path == "" || path == "." {
		return nil
	}
	return os.MkdirAll(path, 0o755)
}
