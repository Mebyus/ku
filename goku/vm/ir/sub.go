package ir

import "github.com/mebyus/ku/goku/vm/opc"

// TestVal compare register and value.
type TestVal struct {
	nodeAtom

	Dest opc.Register
	Val  uint64
}
