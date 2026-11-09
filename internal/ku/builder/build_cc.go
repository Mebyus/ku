package builder

import "github.com/mebyus/ku/internal/ku/builder/kub"

// CompileFileForTestGenC compiles a single source file specified by
// its path into native test executable via genc + cc invocation.
//
// Returns path to resulting executable.
//
// Provided directory will be used for storing temporary intermediate
// files and resulting executable.
//
// Used for C backend (genc) integration tests.
func CompileFileForTestGenC(lang, dir string, path string) (string, error) {
	r := Build(&Config{
		LangDir: lang,
		File:    path,
		OutDir:  dir,
		Build: kub.Build{
			Goal: kub.IntExe,
		},
	})
	if r.Error != nil {
		return "", r.Error
	}

	return r.OutPath, nil
}
