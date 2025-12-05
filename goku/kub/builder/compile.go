package builder

import (
	"github.com/mebyus/ku/goku/compiler/enums/bk"
	"github.com/mebyus/ku/goku/compiler/srcmap"
)

type CompileConfig struct {
	IncludeDirs []string

	// Optional. Specifies directory to be used as root for storing test executable
	// as well as various intermediate files.
	//
	// Default value is used if this field is empty.
	OutputRootDir string

	// Optional. Path to output executable or object file.
	//
	// When this field is empty default value (based on OutputRootDir) will be used.
	OutputPath string

	// Required.
	BuildKind bk.Kind
}

// CompileTexts compile one or more Ku or C source text into one object file.
func CompileTexts(c *CompileConfig, texts []*srcmap.Text) error {
	if len(texts) == 0 {
		panic("no texts")
	}
	return nil
}
