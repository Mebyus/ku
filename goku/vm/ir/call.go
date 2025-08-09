package ir

import "github.com/mebyus/ku/goku/vm/opc"

// FunName contains function name in integer form.
//
// Directly corresponds to function entry index inside list of
// all program functions.
type FunName uint32

// CallFun represents call to address encoded in instruction.
//
// This call is dispatched statically.
type CallFun struct {
	nodeAtom

	Fun FunName
}

// CallReg represents call to address stored in register instruction.
//
// This call is dispatched dynamically.
type CallReg struct {
	nodeAtom

	Reg opc.Register
}
