package tokens

// literal maps Kind to token static literal or code string.
var literal = [...]string{
	empty: "<nil>",

	EOF: "EOF",

	// Operators/punctuators

	Colon:     ":",
	Semicolon: ";",
	Period:    ".",
	Comma:     ",",

	LeftCurly:   "{",
	RightCurly:  "}",
	LeftSquare:  "[",
	RightSquare: "]",
	// LeftParen:      "(",
	// RightParen:     ")",

	// Keywords

	Fun:   "#fun",
	Data:  "#data",
	Entry: "#entry",

	// Special literals

	// Non static literals

	Illegal:    "ILG",
	Word:       "WORD",
	String:     "STR",
	Rune:       "RUNE",
	BinInteger: "INT.BIN",
	OctInteger: "INT.OCT",
	DecInteger: "INT.DEC",
	HexInteger: "INT.HEX",
	DecFloat:   "FLT.DEC",

	// Comments

	LineComment:  "COM.LINE",
	BlockComment: "COM.BLOCK",
}

func (k Kind) String() string {
	return literal[k]
}
