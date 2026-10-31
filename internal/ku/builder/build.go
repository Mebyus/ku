package builder

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mebyus/ku/goku/compiler/cc"
	"github.com/mebyus/ku/goku/compiler/enums/bk"
	"github.com/mebyus/ku/internal/ku/builder/kub"
	"github.com/mebyus/ku/internal/ku/genc"
	"github.com/mebyus/ku/internal/ku/stg"
	"github.com/mebyus/ku/internal/ku/sx"
)

type Config struct {
	Build kub.Build

	// System path of language root directory.
	// Used for resolving units from standard library.
	//
	// Default value will be detected based on compiler executable location.
	LangDir string

	// System path of program/project source root directory.
	// Used for resolving units local to project.
	//
	// Default value will be detected based on current working directory.
	RootDir string

	// Local unit path or system path (must start with "./src/") of what to build. Always a local unit.
	//
	// Required for any build.
	Unit string

	// System path of starting file.
	//
	// Alternative for Unit field.
	//
	// If not empty build starts not from a unit directory,
	// but from a single file which is treated as separate unit.
	//
	// Used for integration tests.
	File string

	// Unit path of build entry (which unit build starts from).
	//
	// Filled based on Unit field.
	unit sx.Path

	// System path to directory which will be used as output directory for various
	// intermediate files and/or build result artifacts.
	OutDir string

	// Path to output translated C code, object file or executable produced as build result.
	// Note that different build types produce different kind of artifacts.
	//
	// Default value will be used if empty.
	OutPath string

	// Build name. Used to generate names for default output name
	// and various intermediate files.
	//
	// Default value will be used if empty (based on unit name).
	Name string
}

// Result describes build result.
type Result struct {
	// May be empty if build resulted in error.
	OutPath string

	// Not nil if build was interrupted by error.
	Error error
}

func Build(config *Config) *Result {
	const debug = false
	start := time.Now()
	r := &Result{}

	err := config.setDefaults()
	if err != nil {
		r.Error = err
		return r
	}

	w := walker{
		std: filepath.Join(config.LangDir, "src", "std"),
		loc: filepath.Join(config.RootDir, "src"),
	}
	w.init(witem{
		path: config.unit,
		file: config.File,
	})
	w.walk()

	n := 0 // total number of errors

	n += len(w.errors)
	for _, e := range w.errors {
		w.pool.FormatError(os.Stderr, e)
	}

	for _, u := range w.units {
		n += len(u.errors)
		for _, e := range u.errors {
			// TODO: refactor into sx.Error
			pos := w.pool.FormatPin(e.Pin)
			fmt.Fprintf(os.Stderr, "%s: %s\n", pos, e.Short)
		}
	}

	if n != 0 && len(w.units) == 0 {
		// failed to load units due to previous errors
		os.Exit(1)
	}

	tower, cycle := rankUnits(w.units, w.roots)
	if len(cycle) != 0 {
		// this error breaks program no import cycle assumption
		// it is fatal to our model of compilation
		fmt.Fprintf(os.Stderr, "import cycle found:\n")
		for _, id := range cycle {
			ipath := w.units[id].path
			fmt.Fprintf(os.Stderr, "    [%d] %s\n", id, ipath)
		}
		os.Exit(1)
	}

	var prog stg.Program
	pool := stg.NewPool(config.Build, w.pool)
	for _, cohort := range tower.Cohorts {
		for _, id := range cohort.Symbols {
			u := w.units[id]

			typer := pool.Get()
			unit := &stg.Unit{
				Path: u.path,
				Dir:  u.dir,
				Name: u.name,
			}
			typer.Do(unit, u.texts)
			pool.Put(typer)

			n += len(unit.Errors)
			for _, e := range unit.Errors {
				pos := w.pool.DecodePin(e.Pin)
				fmt.Fprintf(os.Stderr, "%s: %s\n", pos, e.Short)
			}

			prog.Units = append(prog.Units, unit)
		}
	}
	stg.AssignLinkNames(prog.Units)

	if debug {
		fmt.Printf("uwalk: %s\n", time.Since(start))
	}

	if n != 0 {
		os.Exit(1)
	}

	// TODO: should we pass Common to pool by pointer from Program instead?
	prog.Common = pool.Common

	progName := filepath.Join(config.OutDir, "genc", config.Name+".c")

	start = time.Now()
	err = genProg(progName, &prog)
	if err != nil {
		fmt.Fprintf(os.Stderr, "genc: %s\n", err)
		os.Exit(1)
	}
	if debug {
		fmt.Printf("genc:  %s\n", time.Since(start))
	}

	{
		dir := filepath.Dir(config.OutPath)
		err := os.MkdirAll(dir, 0o755)
		if err != nil {
			r.Error = err
			return r
		}
	}

	start = time.Now()
	err = cc.CompileObj(config.OutPath, progName, bk.Debug)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cc: %s\n", err)
		os.Exit(1)
	}
	if debug {
		fmt.Printf("cc:    %s\n", time.Since(start))
	}

	return r
}

func (c *Config) setDefaults() error {
	if c.Unit == "" && c.File == "" {
		panic("empty build start")
	}

	if c.Unit != "" {
		if strings.HasPrefix(c.Unit, "./src/") {
			p := filepath.Clean(strings.TrimPrefix(c.Unit, "./src/"))
			if p == "" || strings.HasPrefix(p, ".") || strings.HasPrefix(p, "/") {
				return fmt.Errorf("invalid unit path \"%s\"", c.Unit)
			}
			c.unit = sx.MakePath(sx.Loc, p)
		} else {
			c.unit = sx.MakePath(sx.Loc, c.Unit)
		}
	}
	// TODO: validate resulting unit path

	if c.LangDir == "" {
		dir, err := GetLangDir()
		if err != nil {
			return err
		}
		c.LangDir = dir
	}

	if c.RootDir == "" {
		dir, err := os.Getwd()
		if err != nil {
			return err
		}
		c.RootDir = dir
	}

	if c.OutDir == "" {
		c.OutDir = filepath.Join(c.RootDir, outDirName)
	}

	if c.Name == "" {
		var name string
		switch {
		case c.File != "":
			name = strings.TrimSuffix(filepath.Base(c.File), ".ku")
		case c.Unit != "":
			name = filepath.Base(c.Unit)
		}
		if name == "" {
			panic("empty unit name")
		}
		c.Name = name
	}

	if c.OutPath == "" {
		c.OutPath = filepath.Join(c.OutDir, "obj", c.Name+".o")
	}

	return nil
}

const outDirName = ".kub"

// GetOutDir returns path to default build output directory.
func GetOutDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, outDirName), nil
}

// GetLangDir determine Ku root directory based on currently running executable.
func GetLangDir() (string, error) {
	path, err := os.Executable()
	if err != nil {
		return "", err
	}
	if path == "" || path == "." || path == "/" {
		panic("invalid executable path")
	}

	// exec_path = lang_root/bin/ku
	return filepath.Dir(filepath.Dir(path)), nil
}

// genProg generates C code into specified output file.
func genProg(out string, prog *stg.Program) error {
	dir := filepath.Dir(out)
	err := os.MkdirAll(dir, 0o755)
	if err != nil {
		return err
	}

	file, err := os.Create(out)
	if err != nil {
		return err
	}
	defer file.Close()

	return genc.Gen(file, prog)
}
