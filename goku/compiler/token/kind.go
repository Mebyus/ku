package token

type Kind uint32

const (
	// Zero value of Kind. Should not be used explicitly.
	//
	// Mostly a trick to detect places where Kind is left unspecified.
	empty Kind = iota

	// Special tokens
	EOF

	// Operators and/or punctuators

	Underscore // _

	Address     // .&
	Deref       // .*
	DerefIndex  // .[
	BagSelect   // .(
	DerefSelect // .*.

	Plus      // +
	Minus     // -
	Asterisk  // *
	Slash     // /
	Percent   // %
	Ampersand // &
	Quest     // ?

	Pipe       // |
	Caret      // ^
	LeftShift  // <<
	RightShift // >>
	BitAndNot  // &^

	Semicolon // ;
	Period    // .
	Colon     // :
	Comma     // ,
	Ellipsis  // ...

	Equal          // ==
	NotEqual       // !=
	LessOrEqual    // <=
	GreaterOrEqual // >=
	LeftAngle      // <
	RightAngle     // >
	Not            // !

	Assign    // =
	Walrus    // :=
	AddAssign // +=
	SubAssign // -=
	MulAssign // *=
	DivAssign // /=
	RemAssign // %=
	AndAssign // &=
	OrAssign  // |=

	LeftShiftAssign  // <<=
	RightShiftAssign // >>=

	And        // &&
	Or         // ||
	LeftArrow  // <-
	RightArrow // =>

	// Brackets

	LeftCurly   // {
	RightCurly  // }
	LeftSquare  // [
	RightSquare // ]
	LeftParen   // (
	RightParen  // )

	HashCurly  // #{
	HashSquare // #[

	Tweak        // .{
	Insist       // .!
	Chain        // .?
	Chunk        // []
	AutoLen      // [_]
	ArrayPointer // [*]
	ArrayRef     // [&]
	CapBuffer    // [^]

	Nillable      // ?|
	NillableChunk // [?]

	// Keywords

	If
	Else
	In
	For
	Jump

	Defer
	Fun
	Import
	Test

	Gen
	Bag

	Ret

	Struct
	Union

	Const
	Type
	Let
	Var

	Pub
	Unsafe

	Must
	Panic
	Cast // cast
	Tint // tint - truncate (cast with storage size change) integer

	// Special literals

	Nil
	True
	False
	Any // designator to use as *any (void* analog)

	StaticMust // #must
	StaticIf   // #if

	Debug  // #debug
	Build  // #build
	Never  // #never
	Stub   // #stub
	Size   // #size - query type size
	Lookup // #lookup

	TypeId  // #typeid
	ErrorId // #error
	Enum    // #enum

	Check // #check
	Len   // #len

	LabelNext // @.next
	LabelOut  // @.out

	DirName    // #name
	DirInclude // #include
	DirDefine  // #define
	DirLink    // #link

	staticLiteralEnd
)

const (
	// Common tokens with baselex package.
	// Order of common tokens must the same as in baselex.

	Illegal Kind = staticLiteralEnd + iota // any byte sequence unknown to lexer

	// Identifiers and basic type literals
	Word       // myvar, main, Line, print
	BinInteger // 0b1101100001
	OctInteger // 0o43671
	DecInteger // 5367, 43432, 1000097
	HexInteger // 0x43da1
	DecFloat   // 123.45
	Rune       // 'a', '\t', 'p'
	String     // "abc", "", "\t\n  42Hello\n"
	CString    // c"Hello" - C string literal (that is translated to naked .rodata pointer during compilation), terminated with zero byte.

	// Custom tokens

	RawString  // #"raw string literal"
	FillString // "string with ${10 + 1} interpolated ${a - b} expressions"
	Macro      // #.MACRO_NAME
	Env        // #:ENV_NAME

	// Comments
	LineComment  // Line comment starts with //
	BlockComment // Comment inside /* comment */ block

	maxKind
)

// FromBaseKind transforms baselex token kind into Kind.
func FromBaseKind(k uint32) Kind {
	return Kind(k) + staticLiteralEnd
}
