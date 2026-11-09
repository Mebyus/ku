// Unsigned integers of fixed size.
typedef unsigned char      u8;
typedef unsigned short int u16;
typedef unsigned int       u32;
typedef unsigned long int  u64;
typedef __uint128_t        u128;

// Signed integers of fixed size.
typedef signed char      s8;
typedef signed short int s16;
typedef signed int       s32;
typedef signed long int  s64;
typedef __int128_t       s128;

// Target architecture dependant integer types.
typedef u64 uint;
typedef s64 sint;

typedef float      f32;
typedef double     f64;
typedef __float128 f128;

// This should only be used with pointer types.
#define nil 0

// Do a typecast.
#define cast(T, x) (T)(x)

// Span of raw bytes or string with known byte length.
//
// May also be a string that contains C string pointer
// with known byte length attached to it.
typedef struct {
    // Array pointer to raw bytes.
    u8* ptr;

    // Number of bytes available in {ptr}.
    uint num;
} s_u8, str, c_string;

// construct a runtime string from its bytes array pointer and number of bytes
static str
make_str(u8* ptr, uint num) {
    str s = {};
    if (num == 0) {
        return s;
    }

    s.ptr = ptr;
    s.num = num;
    return s;
}

static uint
unsafe_get_c_string_byte_length(u8* ptr) {
	uint i = 0;
	while (ptr[i] != 0) {
		i += 1;
	}
	return i;
}

static c_string
unsafe_make_c_string_from_ptr(u8* ptr) {
	uint num = unsafe_get_c_string_byte_length(ptr);
	c_string s = {};
	s.ptr = ptr;
	s.num = num;
	return s;
}

// Span of strings.
typedef struct {
    // Array pointer to span strings.
    str* ptr;

    // Number of elements available in span.
    uint num;
} s_str;

static s_str
make_span_str(str* ptr, uint num) {
    s_str s = {};
    if (num == 0) {
        return s;
    }

    s.ptr = ptr;
    s.num = num;
    return s;
}
