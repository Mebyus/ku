package sx

import (
	"hash/fnv"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Pool loads and stores source files. Stored files can be accessed as Text.
// Can decode Pin to FilePos if it comes from one of stored files.
type Pool struct {
	// List of all stored texts in order they were loaded.
	list []*Text

	m map[ /* path */ string]*Text
}

func New() *Pool {
	return &Pool{m: make(map[string]*Text)}
}

func (p *Pool) DecodePin(pin Pin) FilePos {
	pos := pin.Pos()
	text := p.get(pos.Text)
	if text == nil {
		return FilePos{}
	}
	if pos.Offset > uint32(len(text.Data)) {
		return FilePos{}
	}
	return FilePos{
		Path: text.Path,
		Pos:  FindTextPos(text.Data, pos.Offset),
	}
}

// Load loads a file by given path and stores it into internal cache.
// Returns Text created from loaded file.
// If file was already loaded previously, then cached version is used.
//
// Path argument should be cleaned by caller for consistency.
func (p *Pool) Load(path string) (*Text, error) {
	if path == "" || path == "." {
		panic("empty or invalid path")
	}
	text, ok := p.m[path]
	if ok {
		return text, nil
	}

	text, err := loadTextFromFile(path)
	if err != nil {
		return nil, err
	}

	p.add(path, text)
	return text, nil
}

func (p *Pool) add(path string, text *Text) {
	id := uint32(len(p.list)) + 1
	text.ID = id

	p.list = append(p.list, text)
	p.m[path] = text
}

// get returns stored Text by its id.
// Returns nil if specified Text not found in cache.
func (p *Pool) get(id uint32) *Text {
	if id == 0 || id > uint32(len(p.list)) {
		return nil
	}
	return p.list[id-1]
}

// does not set Text.ID
func loadTextFromFile(path string) (*Text, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return nil, err
	}

	h := fnv.New64a()
	size := info.Size()
	if size <= 0 {
		return &Text{
			Path: path,
			Ext:  filepath.Ext(path),
			Hash: h.Sum64(),
		}, nil
	}

	// TODO: do we still need this trick when using strings.Builder?
	size += 1 // one byte for final read at EOF

	var g strings.Builder
	g.Grow(int(size))
	_, err = io.Copy(&g, io.TeeReader(f, h))
	if err != nil {
		return nil, err
	}

	return &Text{
		Data: g.String(),
		Path: path,
		Ext:  filepath.Ext(path),
		Hash: h.Sum64(),
	}, nil
}
