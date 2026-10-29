/*
Defines various constants and structures which describe
build target and conditions.
*/
package kub

import (
	"errors"
	"fmt"
	"strings"
)

// Type determines build type. It controls level of optimizations,
// strictness of compiler diagnostics and sets some compile-time
// constants which can be accessed in user code for conditional
// compilation and compile-time decisions.
type Type uint8

const (
	// Debug-friendly optimizations + debug information in binaries + safety checks
	Debug Type = iota + 1

	// Moderate-level optimizations + debug information in binaries + safety checks
	Test

	// Most optimizations enabled + safety checks
	Safe

	// All optimizations enabled + disabled safety checks
	Fast

	_num
)

var buildTypeText = [...]string{
	0: "<nil>",

	Debug: "debug",
	Test:  "test",
	Safe:  "safe",
	Fast:  "fast",
}

func (t Type) String() string {
	return buildTypeText[t]
}

func (t Type) Valid() error {
	if t == 0 {
		return errors.New("empty build type")
	}
	if t >= _num {
		return fmt.Errorf("invalid build type (=%d)", t)
	}
	return nil
}

func Parse(s string) (Type, error) {
	s = strings.TrimSpace(s)

	var t Type
	switch s {
	case "":
		return 0, fmt.Errorf("empty build type")
	case "debug":
		t = Debug
	case "test":
		t = Test
	case "safe":
		t = Safe
	case "fast":
		t = Fast
	default:
		return 0, fmt.Errorf("unknown \"%s\" build type", s)
	}

	if t == 0 {
		panic("empty build type")
	}

	return t, nil
}
