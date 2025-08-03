package builder

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/mebyus/ku/goku/compiler/cc"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/enums/bk"
	"github.com/mebyus/ku/goku/compiler/enums/bm"
	"github.com/mebyus/ku/goku/compiler/srcmap"
	"github.com/mebyus/ku/goku/compiler/srcmap/origin"
	"github.com/mebyus/ku/goku/kub/eval"
	"github.com/mebyus/ku/goku/kub/parser"
)

const defaultOutputRootDir = ".kubout"

// Config contains settings common for various module building tasks.
type Config struct {
	// Required.
	//
	// Specifies which module in package should be tested.
	//
	// Leads to error if there is no module with such name inside the package.
	ModuleName string

	// Required.
	PackageFilePath string

	// Optional. Specifies directory to be used as root for storing test executable
	// as well as various intermediate files.
	//
	// Default value is used if this field is empty.
	OutputRootDir string

	// Required.
	BuildKind bk.Kind
}

type TestModuleConfig struct {
	Config

	// Optional. When this field is not empty, specifies test to be selected
	// from all gathered tests. All other tests are dropped and only selected one
	// gets included into test executable.
	//
	// Leads to error if there is no test with such name inside the module.
	//
	// Useful when user needs to debug chosen test by launching executable inside debugger.
	TestName string
}

type BuildAndRunModuleConfig struct {
	Config
}

type BuildModuleConfig struct {
	Config

	// Optional. Path to output executable or object file.
	//
	// When this field is empty default value (based on OutputRootDir) will be used.
	OutputPath string
}

func (c *Config) checkAndSetDefaults() error {
	err := c.BuildKind.Valid()
	if err != nil {
		return err
	}
	if c.PackageFilePath == "" {
		return errors.New("empty package file path")
	}
	if c.ModuleName == "" {
		return errors.New("empty module name")
	}

	if c.OutputRootDir == "" {
		c.OutputRootDir = defaultOutputRootDir
	}
	return nil
}

func (c *Config) getCodegenOutPath(name string) string {
	return filepath.Join(c.OutputRootDir, "genc", name+".kubgen.c")
}

func (c *Config) getExecutablePath(name string) string {
	return filepath.Join(c.OutputRootDir, "exe", name)
}

func (c *Config) getObjectPath(name string) string {
	return filepath.Join(c.OutputRootDir, "obj", name+".o")
}

// Returns path to file where generated C code for a given module name should be stored.
func (c *Config) getTestCodegenOutPath(name string) string {
	return filepath.Join(c.OutputRootDir, "genc", name+".test.kubgen.c")
}

// Returns path to file where produced test executable for a given module name should be stored.
func (c *Config) getTestExecutablePath(name string) string {
	return filepath.Join(c.OutputRootDir, "test", name)
}

func BuildModule(config *BuildModuleConfig) error {
	err := config.checkAndSetDefaults()
	if err != nil {
		panic(err)
	}

	pool := srcmap.New()
	pkg, err := loadPackageFromFile(pool, config.PackageFilePath)
	if err != nil {
		return err
	}

	name := config.ModuleName
	mod := getPackageModule(pkg, name)
	if mod == nil {
		return fmt.Errorf("package has no module \"%s\"", name)
	}
	if mod.Main == "" {
		return buildObjectFromUnits(config, pool, ResolveConfig{
			RootDir: pkg.RootDir,
			MainDir: pkg.MainDir,
			UnitDir: pkg.UnitDir,
		}, mod.Units)
	}

	codegenOutPath := config.getCodegenOutPath(name)
	err = mkdirForFile(codegenOutPath)
	if err != nil {
		return err
	}

	err = genUnitsToFile(&BuildConfig{
		Pool:      pool,
		BuildKind: config.BuildKind,
		BuildMode: bm.Exe,
		ResolveConfig: ResolveConfig{
			RootDir: pkg.RootDir,
			MainDir: pkg.MainDir,
			UnitDir: pkg.UnitDir,
		},
	},
		codegenOutPath, []srcmap.QueueItem{{
			Path: origin.Path{
				Import: mod.Main,
				Origin: origin.Main,
			},
		}},
	)
	if err != nil {
		return err
	}

	var exePath string
	if config.OutputPath != "" {
		exePath = config.OutputPath
	} else {
		exePath = config.getExecutablePath(name)
	}
	err = mkdirForFile(exePath)
	if err != nil {
		return err
	}

	return cc.CompileExe(exePath, codegenOutPath, config.BuildKind)
}

func buildObjectFromUnits(config *BuildModuleConfig, pool *srcmap.Pool, resolve ResolveConfig, units []string) error {
	name := config.ModuleName
	codegenOutPath := config.getCodegenOutPath(name)
	err := mkdirForFile(codegenOutPath)
	if err != nil {
		return err
	}

	err = genUnitsToFile(&BuildConfig{
		Pool:          pool,
		BuildKind:     config.BuildKind,
		BuildMode:     bm.Obj,
		ResolveConfig: resolve,
	},
		codegenOutPath, makeQueueItems(units))
	if err != nil {
		return err
	}

	var exePath string
	if config.OutputPath != "" {
		exePath = config.OutputPath
	} else {
		exePath = config.getObjectPath(name)
	}
	err = mkdirForFile(exePath)
	if err != nil {
		return err
	}

	return cc.CompileObj(exePath, codegenOutPath, config.BuildKind)
}

func BuildAndRunModule(config *BuildAndRunModuleConfig) error {
	err := config.checkAndSetDefaults()
	if err != nil {
		panic(err)
	}

	pool := srcmap.New()
	pkg, err := loadPackageFromFile(pool, config.PackageFilePath)
	if err != nil {
		return err
	}

	name := config.ModuleName
	mod := getPackageModule(pkg, name)
	if mod == nil {
		return fmt.Errorf("package has no module \"%s\"", name)
	}
	if mod.Main == "" {
		return fmt.Errorf("module \"%s\" does not specify main unit", name)
	}

	codegenOutPath := config.getCodegenOutPath(name)
	err = mkdirForFile(codegenOutPath)
	if err != nil {
		return err
	}

	err = genUnitsToFile(&BuildConfig{
		Pool:      pool,
		BuildKind: config.BuildKind,
		BuildMode: bm.Exe,
		ResolveConfig: ResolveConfig{
			RootDir: pkg.RootDir,
			MainDir: pkg.MainDir,
			UnitDir: pkg.UnitDir,
		},
	},
		codegenOutPath, []srcmap.QueueItem{{
			Path: origin.Path{
				Import: mod.Main,
				Origin: origin.Main,
			},
		}},
	)
	if err != nil {
		return err
	}

	exePath := config.getExecutablePath(name)
	err = mkdirForFile(exePath)
	if err != nil {
		return err
	}

	err = cc.CompileExe(exePath, codegenOutPath, config.BuildKind)
	if err != nil {
		return err
	}

	return runExecutable(exePath)
}

// TestModule builds and runs test executable for a package module.
func TestModule(config *TestModuleConfig) error {
	err := config.checkAndSetDefaults()
	if err != nil {
		panic(err)
	}

	pool := srcmap.New()
	pkg, err := loadPackageFromFile(pool, config.PackageFilePath)
	if err != nil {
		return err
	}

	name := config.ModuleName
	mod := getPackageModule(pkg, name)
	if mod == nil {
		return fmt.Errorf("package has no module \"%s\"", name)
	}

	codegenOutPath := config.getTestCodegenOutPath(name)
	err = mkdirForFile(codegenOutPath)
	if err != nil {
		return err
	}

	err = genUnitsToFile(&BuildConfig{
		Pool:      pool,
		BuildKind: config.BuildKind,
		BuildMode: bm.TestExe,
		ResolveConfig: ResolveConfig{
			RootDir: pkg.RootDir,
			MainDir: pkg.MainDir,
			UnitDir: pkg.UnitDir,
		},
	}, codegenOutPath, makeQueueItems(mod.Units))
	if err != nil {
		return err
	}

	testExePath := config.getTestExecutablePath(name)
	err = mkdirForFile(testExePath)
	if err != nil {
		return err
	}

	err = cc.CompileExe(testExePath, codegenOutPath, config.BuildKind)
	if err != nil {
		return err
	}

	return runTestExecutable(testExePath)
}

func makeQueueItems(ss []string) []srcmap.QueueItem {
	if len(ss) == 0 {
		return nil
	}

	items := make([]srcmap.QueueItem, 0, len(ss))
	for _, s := range ss {
		items = append(items, srcmap.QueueItem{
			Path: origin.Local(s),
			Pin:  0, // TODO: we should bring pins from parsed package file
		})
	}
	return items
}

func getPackageModule(pkg *eval.Package, name string) *eval.Module {
	for _, m := range pkg.Modules {
		if m.Name == name {
			return &m
		}
	}
	return nil
}

func loadPackageFromFile(pool *srcmap.Pool, path string) (*eval.Package, error) {
	text, err := pool.Load(path)
	if err != nil {
		return nil, err
	}

	p := parser.FromText(text)
	astPkg, err := p.Package()
	if err != nil {
		return nil, diag.Format(pool, err.(diag.Error))
	}

	pkg, err := eval.EvalPackage(astPkg)
	if err != nil {
		return nil, diag.Format(pool, err.(diag.Error))
	}
	return pkg, nil
}

func runExecutable(path string) error {
	cmd := exec.Command(path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

const defaultRunTestExecutableTimeout = 30 * time.Second

func runTestExecutable(path string) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultRunTestExecutableTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Creates directory which can house a file with specified path.
func mkdirForFile(path string) error {
	return mkdir(filepath.Dir(path))
}

func mkdir(path string) error {
	if path == "" || path == "." {
		return nil
	}
	return os.MkdirAll(path, 0o755)
}
