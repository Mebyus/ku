/*
Package kvx describes Ku VM executable binary format and
provides functions and data structures to encode and decode it.

Format encodes a Ku VM program. Program consists of 3 segments:
  - Text (sequence of instructions for VM)
  - Data (read-only data loaded by VM)
  - Global (memory to store global variables during execution)

Format stores Text and Data segments as raw binary array of bytes
with offset and size information encoded in its Header.

Global segment only has entry inside the Header, because it is
always initialized to zero, and thus format does not need to
encode it.

Each segment with raw binary is aligned by 8-byte boundary.

All integers are stored in little endian.

Byte layout:

	0: [XX XX XX XX] 4 // Magic
	4: [XX XX XX XX] 4 // Version
	8: [XX XX XX XX] 4 // Entry point offset

	12: [XX XX XX XX XX XX XX XX] 8 // File offset of Text segment
	20: [XX XX XX XX]             4 // Text segment size in bytes
	24: [XX XX XX XX]             4 // Text segment flags
*/
package kvx
