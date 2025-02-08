package source

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

const maxFileSize = 1 << 26

// DirScanParams
type DirScanParams struct {
	// Default value will be used if this field equals 0.
	MaxFileSize uint32

	IncludeTestFiles bool
}

// LoadDir accepts path to directory and optional scan parameters to load all suitable
// source files from the directory. Only the specified directory is scanned, this
// function does no recursive walking. Empty files are always skipped. Ignores
// symbolic links.
//
// Second argument can be nil. In that case default scan parameters are used.
func (p *Pool) LoadDir(dir string, params *DirScanParams) ([]*Text, error) {
	if dir == "" {
		panic("empty unit directory path")
	}

	if params == nil {
		params = &DirScanParams{}
	}
	if params.MaxFileSize == 0 {
		params.MaxFileSize = maxFileSize
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	if len(entries) == 0 {
		return nil, fmt.Errorf("directory \"%s\" is empty", dir)
	}

	var names []string // contains names of selected files
	for _, entry := range entries {
		if !entry.Type().IsRegular() {
			continue
		}

		name := entry.Name()
		if !strings.HasSuffix(name, ".ku") {
			continue
		}
		if !params.IncludeTestFiles && strings.HasSuffix(name, ".test.ku") {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			return nil, err
		}

		if !info.Mode().IsRegular() {
			continue
		}

		size := uint64(info.Size())
		if size == 0 {
			continue
		}
		if size > uint64(params.MaxFileSize) {
			return nil, fmt.Errorf("file \"%s\" is larger (%d bytes) than max allowed size", name, size)
		}

		names = append(names, name)
	}
	if len(names) == 0 {
		return nil, fmt.Errorf("directory \"%s\" does not contain suitable ku source files", dir)
	}
	slices.Sort(names)

	texts := make([]*Text, 0, len(names))
	for _, name := range names {
		text, err := p.Load(filepath.Join(dir, name))
		if err != nil {
			return nil, err
		}
		texts = append(texts, text)
	}
	return texts, nil
}
