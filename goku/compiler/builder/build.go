package builder

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/enums/bk"
	"github.com/mebyus/ku/goku/compiler/srcmap/origin"
)

// Phase specifies build or compilation phase.
type Phase uint8

const (
	// Intermediate code (we use C for now) generation phase.
	PhaseGenC = iota + 1

	// Object file compilation phase.
	PhaseObj

	// Executable linkage phase. Suitable only for targets with entrypoint.
	PhaseExe

	// Test executable output.
	PhaseTest
)

func (p Phase) String() string {
	switch p {
	case 0:
		panic("empty value")
	case PhaseGenC:
		return "genc"
	case PhaseObj:
		return "obj"
	case PhaseExe:
		return "exe"
	case PhaseTest:
		return "test"
	default:
		panic(fmt.Sprintf("unexpected phase (=%d)", p))
	}
}

// Suffix produce a suitable suffix for output name of this phase.
func (p Phase) Suffix() string {
	switch p {
	case 0:
		panic("empty value")
	case PhaseGenC:
		return ".ku.gen.c"
	case PhaseObj:
		return ".o"
	case PhaseExe:
		return ""
	case PhaseTest:
		return ""
	default:
		panic(fmt.Sprintf("unexpected phase (=%d)", p))
	}
}

// Config specifies how to build a unit.
type Config struct {
	// Root directory which contains Ku (cli binary + standard library).
	// If left empty builder will try to determine this automatically.
	RootDir string

	// Path directory containing source code of local units.
	// Default value "src" will be used if empty.
	SourceDir string

	// Directory which will be used for various intermediate generated files.
	// Default value ".kub" will be used if empty.
	GenDir string

	// Path to output translated C code, object file or executable.
	// Default value will be used if empty.
	OutPath string

	// Unit path of what to build. Always a local unit.
	Unit string

	// Build will stop after this phase is complete.
	// Result of this phase will be placed at output path.
	Phase Phase

	BuildKind bk.Kind
}

// SetDefaults set default values for empty fields and check provided values.
func (c *Config) SetDefaults() error {
	if c.Unit == "" {
		panic("empty unit path")
	}

	if c.RootDir == "" {
		dir, err := GetRootDir()
		if err != nil {
			return err
		}
		c.RootDir = dir
	}

	if c.Phase == 0 {
		c.Phase = PhaseObj
	}
	if c.BuildKind == 0 {
		c.BuildKind = bk.Debug
	}

	if c.GenDir == "" {
		c.GenDir = ".kub"
	}
	if c.SourceDir == "" {
		c.SourceDir = "src"
	}

	if c.OutPath == "" {
		name := filepath.Base(c.Unit)
		c.OutPath = filepath.Join(c.GenDir, c.Phase.String(), name+c.Phase.Suffix())
	}

	return nil
}

// GetRootDir determine Ku root directory based on currently running executable.
func GetRootDir() (string, error) {
	path, err := os.Executable()
	if err != nil {
		return "", err
	}
	if path == "" || path == "." || path == "/" {
		panic("invalid executable path")
	}

	return filepath.Dir(filepath.Dir(path)), nil
}

func Build(c *Config) error {
	err := c.SetDefaults()
	if err != nil {
		return err
	}
	err = build(c)
	if err != nil {
		return err
	}
	return nil
}

func build(c *Config) error {
	bundle, err := Walk(WalkConfig{
		Dir: BaseDirs{
			Std: filepath.Join(c.RootDir, "src/std"),
			Loc: c.SourceDir,
		},
	}, QueueItem{
		Path: origin.Local(c.Unit),
	})
	if err != nil {
		return diag.Format(bundle.Pool, err)
	}

	err = CompileBundle(bundle)
	if err != nil {
		return diag.Format(bundle.Pool, err)
	}

	return nil
}
