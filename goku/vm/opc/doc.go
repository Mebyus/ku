/*
Package opc describes instruction opcodes and layouts recognized
by VM as well as how these instructions are encoded in text segment.

Each instruction contains 3 parts:

	Opcode - 1 byte
	Layout - 1 byte
	Data   - variable byte length (0 for some instructions)

Opcode identifies instruction family and layout determines the exact
operation and data encoding within that family.

Combination of Opcode + Layout determines Data length with no ambiguity.
Although some instructions have no Data (may be considered as length 0).

In general Data encodes immediate values and registers used by instruction.
Each register is encoded with 1 byte. Immediate values are always encoded
in little endian and may be 4 or 8 bytes in size. Registers are always encoded
before immediate values and destination operand comes before source operand.

Trap instruction is special, because it is meant as a safeguard against
execution of empty text memory (for example due to function address alignment).
Thus Trap ignores Layout and any instruction of the form [00 XX ...] always
interrupts execution regardless of Layout.

Some instructions (usually arithmetic, such as Add or Sub) use generic
Layout which is common for many instructions. Other instructions such as
Jump use specialized Layout values.

List of instructions with generic layouts:

	Add
	Sub
	Mul
	Div

	Or
	And
	Xor
	Shift
*/
package opc
