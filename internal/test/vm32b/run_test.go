package vm32b_test

import (
	"bytes"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/mebyus/ku/internal/ku/builder"
	"github.com/mebyus/ku/internal/vm32b"
)

// TestRun is an end-to-end test of Ku compiler and
// 32-bit bytecode virtual machine.
//
// Test imitates real-world scenario when program is compiled and
// saved into output file and then executed via VM. In more detail
// test consists of the following steps:
//
//  1. Compile program from a Ku source file into bytecode program.
//  2. Encode resulting program in KVX format (as if saving it into file).
//  3. Decode program from KVX format.
//  4. Execute program inside VM with given arguments.
//  5. Compare VM output with expected test result.
//
// Each program should be executed multiple times with different
// arguments (expected results change based on given arguments).
//
// This test is positive. All programs used here must compile
// and execute without errors.
//
// Test program sources are stored in "./src" directory.
func TestRun(t *testing.T) {
	bundles := []Bundle{
		{
			// empty start function
			src: "0000_empty_start.ku",
			tests: []TestCase{
				{
					args: nil,
					want: nil,
				},
			},
		},
		{
			// start returns a constant value
			src: "0001_return_const.ku",
			tests: []TestCase{
				{
					args: nil,
					want: []uint32{18},
				},
			},
		},
		{
			// start returns a function argument
			src: "0002_return_arg.ku",
			tests: []TestCase{
				{
					args: []uint32{7},
					want: []uint32{7},
				},
				{
					args: []uint32{36},
					want: []uint32{36},
				},
			},
		},
		{
			// start returns a sum of two arguments
			src: "0003_return_add.ku",
			tests: []TestCase{
				{
					args: []uint32{1, 2},
					want: []uint32{3},
				},
				{
					args: []uint32{7, 9},
					want: []uint32{16},
				},
				{
					args: []uint32{0, 43},
					want: []uint32{43},
				},
			},
		},
		{
			// start returns a local variable
			src: "0004_var_local.ku",
			tests: []TestCase{
				{
					args: nil,
					want: []uint32{96},
				},
			},
		},
		{
			// start returns a mutated local variable
			src: "0005_var_change.ku",
			tests: []TestCase{
				{
					args: []uint32{0},
					want: []uint32{3},
				},
				{
					args: []uint32{1},
					want: []uint32{5},
				},
				{
					args: []uint32{10},
					want: []uint32{23},
				},
			},
		},
		{
			// start returns n-th fibonacci number
			src: "0006_fib_loop.ku",
			tests: []TestCase{
				{
					args: []uint32{0},
					want: []uint32{0},
				},
				{
					args: []uint32{1},
					want: []uint32{1},
				},
				{
					args: []uint32{2},
					want: []uint32{1},
				},
				{
					args: []uint32{3},
					want: []uint32{2},
				},
				{
					args: []uint32{4},
					want: []uint32{3},
				},
				{
					args: []uint32{5},
					want: []uint32{5},
				},
				{
					args: []uint32{6},
					want: []uint32{8},
				},
				{
					args: []uint32{7},
					want: []uint32{13},
				},
				{
					args: []uint32{8},
					want: []uint32{21},
				},
				{
					args: []uint32{9},
					want: []uint32{34},
				},
				{
					args: []uint32{10},
					want: []uint32{55},
				},
				{
					args: []uint32{11},
					want: []uint32{89},
				},
			},
		},
		{
			// use function call with no args
			src: "0007_call_noargs.ku",
			tests: []TestCase{
				{
					args: []uint32{0},
					want: []uint32{2},
				},
				{
					args: []uint32{10},
					want: []uint32{12},
				},
				{
					args: []uint32{777},
					want: []uint32{779},
				},
			},
		},
		{
			// change local variable via reference
			src: "0008_ref_local.ku",
			tests: []TestCase{
				{
					args: []uint32{0},
					want: []uint32{8},
				},
				{
					args: []uint32{5},
					want: []uint32{13},
				},
			},
		},
		{
			// change local variable via reference
			src: "0009_if_block.ku",
			tests: []TestCase{
				{
					args: []uint32{0},
					want: []uint32{6},
				},
				{
					args: []uint32{4},
					want: []uint32{10},
				},
				{
					args: []uint32{16},
					want: []uint32{22},
				},
				{
					args: []uint32{17},
					want: []uint32{0},
				},
				{
					args: []uint32{20},
					want: []uint32{3},
				},
			},
		},
		{
			// discard expression
			src: "0010_discard.ku",
			tests: []TestCase{
				{
					args: []uint32{5},
					want: nil,
				},
			},
		},
		{
			// use short if-flow (with return)
			src: "0011_if_return.ku",
			tests: []TestCase{
				{
					args: []uint32{0},
					want: []uint32{1},
				},
				{
					args: []uint32{5},
					want: []uint32{6},
				},
				{
					args: []uint32{16},
					want: []uint32{3},
				},
				{
					args: []uint32{95814},
					want: []uint32{95815},
				},
			},
		},
	}

	// validate bundles and test cases
	n := 0 // number of validation errors
	set := make(map[string]struct{}, len(bundles))
	for i, b := range bundles {
		name := b.src
		if name == "" {
			t.Errorf("bundle (i=%d) has no source file", i)
			n += 1
		} else {
			_, ok := set[name]
			if ok {
				t.Errorf("duplicate source file \"%s\" in bundles", name)
				n += 1
			}
			set[name] = struct{}{}
		}

		if len(b.tests) == 0 {
			t.Errorf("no test cases in \"%s\" bundle", name)
			n += 1
		}
	}
	if n != 0 {
		return
	}

	var buf bytes.Buffer // stores encoded bytecode after compilation
	var vm vm32b.VM
	for _, b := range bundles {
		name := b.src
		path := filepath.Join("src", name)

		buf.Reset()
		err := builder.CompileEncodeForVM32B(path, &buf)
		if err != nil {
			t.Errorf("compile program \"%s\": %v", name, err)
			continue
		}
		prog, err := vm32b.DecodeFromBytes(buf.Bytes())
		if err != nil {
			t.Errorf("decode program \"%s\": %v", name, err)
			continue
		}

		for _, c := range b.tests {
			exit := vm.Exec(prog, c.args)
			if !exit.IsOK() {
				t.Errorf("program \"%s\": Exec(%s) = RuntimeError (=%d)", name, formatArgs(c.args), exit.Error.Code)
				continue
			}

			got := exit.Ret
			if !slices.Equal(got, c.want) {
				t.Errorf("program \"%s\": Exec(%s) = %v, want %v", name, formatArgs(c.args), got, c.want)
			}
		}
	}
}

func formatArgs(args []uint32) string {
	if len(args) == 0 {
		return ""
	}

	list := make([]string, 0, len(args))
	for _, a := range args {
		list = append(list, strconv.FormatUint(uint64(a), 10))
	}
	return strings.Join(list, ", ")
}

type TestCase struct {
	// program arguments
	args []uint32

	// expected output after execution inside VM
	want []uint32
}

// Bundle is a collection of test cases for a single program.
type Bundle struct {
	tests []TestCase

	// filename of source program
	src string
}
