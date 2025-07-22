package exk

// Kind indicates expression kind.
//
// Some values of Kind are only meant to be used in AST, while others are restriced
// to STG. Each value should document its intended context for usage.
type Kind uint32

const (
	// Zero value of Kind. Should not be used explicitly.
	//
	// Mostly a trick to detect places where Kind is left unspecified.
	empty Kind = iota

	// Operand.
	//
	// AST and STG.
	//
	// Indicates operand which consists of single word token usage.
	//
	// AST examples:
	//	foo + bar // Symbol(foo) + Symbol(bar)
	//	a[i]      // a[Symbol(i)]
	Symbol

	// Operand.
	//
	// AST only.
	DotName

	// Operand.
	//
	// AST and STG.
	//
	// In AST represents unsigned integer literal of any kind (decimal, binary,
	// octal, hex).
	//
	// In STG represents static integer value (possibly size-typed). This includes
	// both positive, negative and zero integers.
	//
	// AST examples:
	//	2 + 5 - 3  // Integer(2) + Integer(5) - Integer(3)
	//	-10        // -Integer(10)
	//	a[1]       // a[Integer(1)]
	//	foo & 0xFF // Symbol(foo) & Integer(0xFF)
	//
	// STG examples:
	//	4    // Integer(4)
	//	-10  // Integer(-10)
	//	a[1] // a[Integer(1)]
	//
	// Reduce transform examples:
	//	2 - 3  => -1 // Integer(2) - Integer(3)  => Integer(-1)
	//	2 << 2 => 8  // Integer(2) << Integer(2) => Integer(8)
	Integer

	// Operand.
	//
	// AST and STG.
	//
	// In AST represents non-negative floating point literal.
	//
	// In STG represents static floating point number (possibly typed). This includes
	// both positive, negative and zero numbers.
	//
	// AST examples:
	//	2.0        // Float(2.0)
	//	2.5 + 0.03 // Float(2.5) + Float(0.03)
	Float

	// Operand
	Rune

	// Operand.
	//
	// AST and STG.
	//
	// In AST represents string literal.
	//
	// In STG represents static string value (possibly typed).
	//
	// Examples:
	//	"hello" // String("hello")
	String

	// Operand.
	//
	// In AST represents true literal.
	//
	// In STG represents static boolean true value.
	True

	// Operand.
	//
	// AST and STG.
	//
	// In AST represents false literal.
	//
	// In STG represents static boolean false value.
	//
	// Reduce transform examples:
	//	5 == 3 => false // Integer(5) == Integer(3) => False
	False

	// Operand.
	//
	// AST and STG.
	//
	// In AST represents nil literal.
	//
	// Examples:
	//	p = nil;
	Nil

	// Special form.
	//
	// AST only.
	//
	// Special literal type. Can only be used as init expression for variables.
	// Means that variable is left with its "dirty" stack state instead of
	// initialization to deterministic value.
	//
	// Usage of this literal as a part of any other expression is illegal.
	// Error will be detected during parsing stage.
	//
	// Examples:
	//	var foo: u32 = ?;
	Dirty

	// Special form.
	//
	// Skips value assignment.
	Blank

	// Examples:
	//	{name: "hello", id: 1}
	Object

	// Examples:
	//	[1, 2, 3]
	//	["hello", "world"]
	//	[]
	List

	// Operand.
	//
	// AST and STG.
	//
	// Represents non-trivial unary expression. Non-trivial means that for example
	// expression:
	//	53 // Integer(53)
	// Cannot be represented as unary expression with no operator.
	//
	// Examples:
	//	-10 // Unary(Minus, Integer(10))
	//	a || !b // Symbol(a) || Unary(Not, Symbol(b))
	Unary

	// Expression.
	//
	// AST and STG.
	//
	// Represents binary expression.
	//
	// Examples:
	//	1 - 2 // Binary(Sub, Integer(1), Integer(2))
	//	a & b // Binary(BitAnd, Symbol(a), Symbol(b))
	Binary

	// Operand.
	//
	// AST only.
	//
	// Represents expression inside parenthesis.
	//
	// Examples:
	//	(1 - 2)         // Paren(Binary(Sub, Integer(1), Integer(2)))
	//	a * (c - a + 2) // Binary(Mul, Symbol(a), Paren(Binary(Add, Binary(Sub, Symbol(c), Symbol(a)), Integer(2))))
	Paren

	// Operand.
	//
	// AST only.
	//
	// Represents chain operand. They always start from word token and
	// contain one or more chain parts.
	//
	// Chain operands in Ku language is a way to healthy restrict various usages
	// of "chaining techniques" from other languages.
	//
	// Examples:
	//	a.foo     // Chain(a, Select(foo))
	//	a[5]      // Chain(a, Index(Integer(5)))
	//	a.foo.bar // Chain(a, Select(foo), Select(bar))
	//
	// Some parts terminate the chain and trying to chain it further is illegal.
	// This restriction greatly simplifies read-write access semantics (lvalue vs
	// rvalue differences). Consider the following expressions:
	//	(ok)      a.foo.&   // Takes the address of a.foo member
	//	(illegal) a.foo().& // Tries to take the address of method call result
	//	(illegal) a.foo.&() // Tries to call the address of a.foo member
	//
	// Parts that terminate a chain always form a separate operand, not a chain.
	// List of chain termintators:
	//	- Chain(...).&    - "address of" operator
	//	- Chain(...)(...) - call
	Chain

	// Chain part.
	//
	// AST only.
	//
	// Represents chain part which selects a member.
	//
	// Examples:
	//	a.foo     // Chain(a, Select(foo))
	//	a.foo.bar // Chain(a, Select(foo), Select(bar))
	//	a[0].foo  // Chain(a, Index(Integer(0)), Select(foo))
	Select

	// Chain part.
	//
	// AST only.
	SelectTest

	// Chain part.
	//
	// AST (TODO: complete usage annotation).
	//
	// Represents chain part which dereferences a pointer.
	// This operation sometimes called indirection.
	//
	// Examples:
	//	a.foo.*     // Chain(a, Select(foo), Deref)
	//	a.foo.bar.* // Chain(a, Select(foo), Select(bar), Deref)
	//	a.*         // Chain(a, Deref)
	//	a[2].*      // Chain(a, Index(Integer(2)), Deref)
	Deref

	// ChainPart
	//
	// Examples:
	//	a.foo.*.len
	DerefSelect

	// Chain part.
	//
	// AST.
	//
	// Represents chain part which does access via index.
	//
	// Examples:
	//	a.foo[1] // Chain(a, Select(foo), Index(1))
	//	a[5]     // Chain(a, Index(Integer(5)))
	//	a[0].foo // Chain(a, Index(Integer(0)), Select(foo))
	//	foo[a]   // Chain(foo, Index(Symbol(a)))
	Index

	// Chain part.
	//
	// AST.
	//
	// Represents chain part which does access via index into array pointer.
	//
	// Examples:
	//	a.[1]        // Chain(a, DerefIndex(Integer(0)))
	//	a.foo.[1]    // Chain(a, Select(foo), DerefIndex(Integer(0)))
	//	a[i].foo.[0] // Chain(a, Index(Symbol(a)), Select(foo), DerefIndex(Integer(0)))
	//	a.[1].foo    // Chain(a, DerefIndex(Integer(0)), Select(foo))
	DerefIndex

	// Operand. Terminates chain.
	//
	// AST (TODO: complete usage annotation).
	//
	// Represents chain terminator that takes address of operand.
	//
	// Examples:
	//	a.&        // Ref(Chain(a))
	//	a.foo.&    // Ref(Chain(a, Select(foo)))
	//	a[5].&     // Ref(Chain(a, Index(Integer(5))))
	//	a[i].foo.& // Ref(Chain(a, Index(Symbol(i)), Select(foo)))
	Ref

	// Operand. Terminates chain.
	Tweak

	// Operand. Terminates chain.
	//
	// AST (TODO: complete usage annotation).
	//
	// Represents chain terminator that does a call to operand.
	//
	// Examples:
	//	foo()       // Call(Chain(foo))
	//	a.foo(4, b) // Call(Chain(a, Select(foo)), Arg(Integer(4)), Arg(Symbol(b)))
	//	a[6].foo(0) // Call(Chain(a, Index(Integer(6)), Select(foo)), Arg(Integer(0)))
	Call

	Slice

	// Special form.
	Pack

	// Operand.
	//
	// AST and STG.
	//
	// Static (known at compile time) expression cast to another type.
	// Types must be compatible for cast to be legal.
	//
	// Function-like construct which corresponds to compiler intrinsic.
	//
	// Examples:
	//	#cast(*any, a)
	//	#cast(Foo, bar)
	Cast

	// Operand
	CheckFlag

	// Operand
	ArrayLen

	// Operand.
	//
	// AST and STG.
	//
	// Static (known at compile time) integer expression cast to integer of another
	// size.
	//
	// Function-like construct which corresponds to compiler intrinsic.
	//
	// Examples:
	//	#tint(u32, a)
	//	#tint(u64, a) + 10
	Tint

	// Operand.
	//
	// AST and STG.
	//
	// Static (known at compile time) query of type size.
	//
	// Function-like construct which corresponds to compiler intrinsic.
	//
	// Examples:
	//	#size(u32)
	//	#size(Foo)
	Size

	// Operand
	//
	// Examples:
	//	#typeid(u32)
	//	#typeid(Foo)
	TypeId

	// Operand
	//
	// Examples:
	//	#error(READ_EOF)
	//	#error(WRITE_EOF)
	ErrorId

	EnumMacro

	maxKind
)

var text = [...]string{
	empty: "<nil>",

	Cast: "cast",
	Tint: "tint",
	Size: "#size",

	TypeId:    "#typeid",
	ErrorId:   "#error",
	EnumMacro: "#enum",

	CheckFlag: "#check",
	ArrayLen:  "#len",

	Unary:  "unary",
	Binary: "binary",
	Paren:  "paren",
	Pack:   "pack",
	List:   "list",
	Object: "object",

	Chain:  "chain",
	Index:  "index",
	Call:   "call",
	Ref:    "ref",
	Deref:  "deref",
	Select: "select",
	Slice:  "slice",
	Tweak:  "tweak",

	DerefIndex:  "deref.index",
	DerefSelect: "deref.select",
	SelectTest:  "select.test",
	DotName:     "dot.name",

	Symbol:  "symbol",
	Integer: "integer",
	Float:   "float",
	Rune:    "rune",
	String:  "string",
	True:    "true",
	False:   "false",
	Nil:     "nil",
	Dirty:   "dirty",
	Blank:   "blank",
}

func (k Kind) String() string {
	return text[k]
}
