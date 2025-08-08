package kvx

// Program holds program information in a form that is suitable
// for VM execution.
type Program struct {
	Text []byte
	Data []byte

	// Offset into program text.
	EntryPoint uint32

	GlobalSize uint32
}

type SegmentHeader struct {
	Offset uint64
	Size   uint32
	Flags  uint32
}

type Header struct {
	Text   SegmentHeader
	Data   SegmentHeader
	Global SegmentHeader

	Version uint32
}

type File struct {
	Header Header

	Program *Program
}

const Magic = "KVX\x00"

const HeaderSize = 4 + 4 + // Magic + Version
	16 + // Text Header
	16 + // Data Header
	16 // Global Header

func NewFile(prog *Program) *File {
	var offset uint64

	offset = alignBy8(HeaderSize)
	textHeader := SegmentHeader{
		Offset: offset,
		Size:   uint32(len(prog.Text)),
	}

	offset = alignBy8(offset + uint64(textHeader.Size))
	dataHeader := SegmentHeader{
		Offset: offset,
		Size:   uint32(len(prog.Data)),
	}

	globalHeader := SegmentHeader{
		Offset: 0, // offset is ignored for this segment
		Size:   prog.GlobalSize,
	}

	return &File{
		Header: Header{
			Text:   textHeader,
			Data:   dataHeader,
			Global: globalHeader,
		},
		Program: prog,
	}
}

func alignBy8(v uint64) uint64 {
	a := v & 0x7
	a = ((^a) + 1) & 0x7
	return v + a
}
