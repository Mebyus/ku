package source

import (
	"fmt"
	"os"
)

// Pool loads and stores source files. Stored files can be accessed as Text.
// Implements PinMap by translating file id from Pin to actual file path.
type Pool struct {
	// List of all stored texts in order they were loaded.
	list []Text

	m map[ /* path */ string]*Text
}

func New() *Pool {
	return &Pool{m: make(map[string]*Text)}
}

func (p *Pool) DecodePin(pin Pin) (FilePos, error) {
	pos := pin.Pos()
	text := p.get(pos.Text)
	if text == nil {
		return FilePos{}, fmt.Errorf("text (id=%d) not found", pos.Text)
	}
	if pos.Offset > uint32(len(text.Data)) {
		return FilePos{}, fmt.Errorf("offset (=%d) out of text (len=%d)", pos.Offset, len(text.Data))
	}
	return FilePos{
		Path: text.Path,
		Pos:  FindTextPos(text.Data, pos.Offset),
	}, nil
}

// Load loads a file by given path and stores it into internal cache.
// Returns Text created from loaded file.
// If file was already loaded previously, then cached version is used.
//
// Path argument should be cleaned by caller for consistency.
func (p *Pool) Load(path string) (*Text, error) {
	t, ok := p.m[path]
	if ok {
		return t, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	i := uint32(len(p.list))
	p.list = append(p.list, Text{
		Data: data,
		Path: path,
		ID:   i + 1,
	})
	t = &p.list[i]
	p.m[path] = t
	return t, nil
}

// get returns stored Text by its id.
// Returns nil if specified Text not found in cache.
func (p *Pool) get(id uint32) *Text {
	if id == 0 {
		panic("zero id")
	}

	if id > uint32(len(p.list)) {
		return nil
	}
	return &p.list[id-1]
}
