package vm32b

import "fmt"

type Prog struct {
	text []inst

	// must containt instruction adderreses (=index into text)
	// in ascending order
	funs []finfo
}

// static information about compiled function
type finfo struct {
	// start instruction index
	addr ia

	// number of arguments required to make a call
	narg uint32

	// number of allocated local variables
	// does not include arguments but does
	// include stack memory gaps
	nloc uint32

	// number of returned (via evs) values
	nret uint32
}

// inst represents code instruction
type inst struct {
	op op

	// not all instructions use this
	//
	//	push  - value pushed on stack
	//	load  - variable name
	//	store - variable name
	//	call  - function name
	//	jump  - instruction address
	//	jz    - instruction address
	arg uint32
}

// holds instruction address which equals instruction
// index inside the whole program
type ia uint32

// holds vm opcode
type op uint32

const (
	// should be zero value
	opTRAP op = iota
	opHALT
	opPRINT
	opCALL
	opRET
	opJUMP  // unconditional jump to target instruction address
	opJZ    // pop EVS and jump to target instruction address if value is 0
	opJNZ   // pop EVS and jump to target instruction address if value is not 0
	opTAKE  // pop value from EVS and place it into stack memory as call argument
	opDRAIN // pop value from EVS and discard it
	opPUSH  // push value stored in local variable onto EVS
	opPOP   // pop value from EVS and store it into local variable
	opPUT   // push constant value onto EVS
	opADD
	opSUB
	opMUL
	opEQUAL
	opLESS

	// opLOAD     // push global variable value onto EVS
	// opSTORE    // pop value from EVS and save it to global variable
)

func PrintProg(prog *Prog) {
	if len(prog.text) == 0 {
		return
	}

	var fun int // current function index
	for i, x := range prog.text {
		if fun < len(prog.funs) && i == int(prog.funs[fun].addr) {
			if fun != 0 {
				fmt.Println()
			}
			fmt.Printf("FUN %d (0x%04X):\n\n", fun, i)
			fun += 1
		}

		fmt.Printf("  %04X:", i)
		switch x.op {
		case opTRAP:
			fmt.Printf("    TRAP\n")
		case opHALT:
			fmt.Printf("    HALT\n")
		case opPRINT:
			fmt.Printf("    PRINT\n")
		case opCALL:
			fmt.Printf("    CALL   #.%d\n", x.arg)
		case opRET:
			fmt.Printf("    RET\n")
		case opJUMP:
			fmt.Printf("    JUMP   0x%X  (%d)\n", x.arg, int(x.arg)-i)
		case opJZ:
			fmt.Printf("    JZ     0x%X  (%d)\n", x.arg, int(x.arg)-i)
		case opJNZ:
			fmt.Printf("    JNZ    0x%X  (%d)\n", x.arg, int(x.arg)-i)
		case opTAKE:
			fmt.Printf("    TAKE\n")
		case opDRAIN:
			fmt.Printf("    DRAIN\n")
		case opPUSH:
			fmt.Printf("    PUSH   [%d]\n", x.arg)
		case opPOP:
			fmt.Printf("    POP    [%d]\n", x.arg)
		case opPUT:
			fmt.Printf("    PUT    0x%X  (%d)\n", x.arg, int32(x.arg))
		case opADD:
			fmt.Printf("    ADD\n")
		case opSUB:
			fmt.Printf("    SUB\n")
		case opMUL:
			fmt.Printf("    MUL\n")
		case opEQUAL:
			fmt.Printf("    EQUAL\n")
		case opLESS:
			fmt.Printf("    LESS\n")
		default:
			panic(fmt.Sprintf("unexpected opcode (=%d)", x.op))
		}
	}
}
