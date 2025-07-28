package bm

import (
	"errors"
	"fmt"
)

// Mode determines build mode.
type Mode uint8

const (
	// Build object file.
	Obj Mode = iota + 1

	// Build executable.
	Exe

	// Build test executable.
	TestExe

	num
)

var buildModeText = [...]string{
	0: "<nil>",

	Obj:     "obj",
	Exe:     "exe",
	TestExe: "exe.test",
}

func (m Mode) String() string {
	return buildModeText[m]
}

func (m Mode) Valid() error {
	if m == 0 {
		return errors.New("empty build mode")
	}
	if m >= num {
		return fmt.Errorf("invalid build mode (=%d)", m)
	}
	return nil
}
