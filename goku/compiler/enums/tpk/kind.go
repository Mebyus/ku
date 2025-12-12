package tpk

// Kind indicates type kind.
type Kind uint8

const (
	// Zero value of Kind. Should not be used explicitly.
	//
	// Mostly a trick to detect places where Kind is left unspecified.
	empty Kind = iota

	// Trivial type has no size or properties. It is formed by language
	// constructs such as:
	//
	//	()        // empty tuple
	//	struct {} // empty struct
	//	[0]int    // array with zero size
	//
	// There is always only one trivial type.
	Trivial

	// Type for any integers (literals, constants and expressions)
	// excluding custom integer types.
	//
	// Flags specify whether this integer type has specified storage
	// size, is signed or unsigned, static or runtime type.
	//
	// Types with fixed bit size which can hold unsigned integers:
	//
	//	u8, u16, u32, u64, u128, uint
	//
	// Types with fixed bit size which can hold signed integers:
	//
	//	s8, s16, s32, s64, s128, sint
	//
	// Static integer type of arbitrary size is specified as:
	//
	//	<int>
	//
	// And this static type has size set to 0.
	Integer

	// Builtin type which can hold fixed number of bytes. Inner structure
	// of this type is identical to chunk of bytes []u8. The difference
	// between these two types is logical, in most cases strings should
	// hold utf-8 encoded text opposed to arbitrary byte sequence in
	// chunk of bytes. However this is not a rule and not enforced by the
	// language in any way. Strings and chunks of bytes can be cast between
	// each other freely with no runtime overhead.
	//
	// Type for any strings (literals, constants and exressions).
	//
	// Static string type is specified as:
	//
	//	<str>
	//
	// And this static type has size set to 0.
	String

	// Builtin type which can hold one of two mutually exclusive values:
	//
	//	- true
	//	- false
	//
	// This type is not a flavor of any integer type. Furthermore boolean
	// flavors cannot be cast into integer flavors via tint operation.
	//
	// Type for any boolean values (literals, constants and exressions).
	//
	// Flags specify whether this boolean type is static or runtime type.
	//
	// Static boolean type is specified as:
	//
	//	<bool>
	//
	// And this static type has size set to 0.
	Boolean

	Float

	AnyPointer

	// Types which are defined as continuous region of memory holding
	// fixed number of elements of another Type
	//
	//	[5]u64 - an array of five u64 numbers
	Array

	Chunk

	// Pointer to another Type. An example of such a Type could be *str
	// (pointer to a string)
	Pointer

	// Pointer to continuous region of memory holding
	// variable number of elements of another type.
	ArrayPointer

	// Types of symbols, created by importing other unit
	//
	//	import std {
	//		math => "math"
	//	}
	//
	// Symbol math in the example above has Type with Kind = Import.
	Import

	// Types which are defined by giving a new name to another type.
	//
	// Created via language construct:
	//
	//	type Name => <OtherType>
	Custom
)

var text = [...]string{
	empty: "<nil>",

	Trivial: "trivial",
	Integer: "int",
	String:  "str",
	Boolean: "bool",
	Float:   "float",

	Array:  "array",
	Chunk:  "chunk",
	Custom: "custom",

	Pointer:      "pointer",
	AnyPointer:   "pointer.any",
	ArrayPointer: "pointer.array",

	Import: "import",
}

func (k Kind) String() string {
	return text[k]
}
