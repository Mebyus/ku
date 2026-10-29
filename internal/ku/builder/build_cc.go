package builder

// CompileFileForTestGenC compiles a single source file specified by
// its path into native test executable via genc + cc invocation.
//
// Returns path to resulting executable.
//
// Provided directory will be used for storing temporary intermediate
// files and resulting executable.
//
// Used for C backend (genc) integration tests.
func CompileFileForTestGenC(dir string, path string) (string, error) {
	r := Build(&Config{
		File:   path,
		OutDir: dir,
	})
	if r.Error != nil {
		return "", r.Error
	}

	return r.OutPath, nil
}
