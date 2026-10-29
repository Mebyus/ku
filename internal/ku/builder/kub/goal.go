package kub

// Goal controls what type of artifact build should produce.
type Goal uint8

const (
	// Build goal will be detected automatically
	// based on starting unit/input.
	Auto Goal = iota

	// Check program for errors.
	//
	// Build as usual up until code generation and
	// then stop.
	//
	// Includes:
	//	- parsing
	//	- import resolution
	//	- typecheck
	Check

	// Build will generate C code.
	GenC

	// Build will produce compiled object file.
	Obj

	// Build will produce native executable.
	Exe

	// Build will produce native executable with
	// test toolchain harness.
	//
	// Used by "ku test" command.
	TestExe

	// Build will produce native executable for
	// integration test.
	//
	// Used by some integration compiler tests.
	IntExe
)
