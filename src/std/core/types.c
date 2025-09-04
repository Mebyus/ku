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

typedef float      f32;
typedef double     f64;
typedef __float128 f128;

// No longer needed in C23
// TODO: make this into conditional compilation for
// different compilers
//
// typedef _Bool bool;

// Represents unicode point.
typedef u32 rune;

// Arch dependant unsigned integer (same size as pointer).
typedef u64 uint;

// Arch dependant signed integer (same size as pointer).
typedef s64 sint;

// Arch dependant pointer in raw integer form. Meant for performing
// arbitrary pointer manipulations.
typedef uint pint;

// This should only be used with pointer types.
#define nil 0

// These should only be used with boolean type.
#define true  1
#define false 0

// Macro for type cast with fancy syntax.
#define cast(T, x) (T)(x)

// Shorthand for C string literal to {str} conversion.
#define ss(s) make_str((u8*)(u8##s), sizeof(u8##s) - 1)

// Alternative for C string literal to {str} conversion when struct initializer is needed.
#define sl(s) { .ptr = (u8*)(u8##s), .len = sizeof(u8##s) - 1 }

// Gives number of elements in non-decayed array.
#define array_len(a) (sizeof(a) / sizeof(*a))

// Produce error code with kind {k} and value {v}
#define ERROR_CODE(k, v) ((k << 32) | v)

#define SOURCE_ORIGIN make_source_origin(make_str((u8*)(__FILE__), sizeof(__FILE__) - 1), __LINE__)

#define panic(s) panic_origin(s, SOURCE_ORIGIN)

_Noreturn static void
panic_trap(void) {
	__builtin_trap();
	__builtin_unreachable();
}

static u32
min_u32(u32 a, u32 b) {
    if (a < b) {
        return a;
    }
    return b;
}

static u32
max_u32(u32 a, u32 b) {
    if (a < b) {
        return b;
    }
    return a;
}

static uint
min_uint(uint a, uint b) {
    if (a < b) {
        return a;
    }
    return b;
}

static uint
max_uint(uint a, uint b) {
    if (a < b) {
        return b;
    }
    return a;
}

/*
Note that {c_string} can be used where {str} or {span_u8} are accepted.
But reverse usage is not valid.
*/
typedef struct {
	// Array pointer to raw bytes data.
	u8* ptr;

	// Number of bytes in array designated by {ptr}.
	//
	// For {c_string} contains byte length of the string (without null byte).
	uint len;
} span_u8, str, c_string;

typedef struct {
	str  file;
	uint line;
} SourceOrigin;

static SourceOrigin
make_source_origin(str file, uint line) {
	SourceOrigin s = {};
	s.file = file;
	s.line = line;
	return s;
}

// Use when you need to explicitly assign empty string.
static const str
empty_str = { .ptr = nil, .len = 0 };

static void
stdout_write(span_u8 s);

static void
stderr_write(span_u8 s);

static void
print(str s) {
	stdout_write(s);
}

static void
panic_origin(str s, SourceOrigin o) {
	stderr_write(s); // TODO: need formatting here
	stderr_write(o.file);
	panic_trap();
}

static span_u8
make_span_u8(u8* ptr, uint len) {
    span_u8 s = {};
	if (len == 0) {
		return s;
	}

    s.ptr = ptr;
    s.len = len;
    return s;
}

static str
make_str(u8* ptr, uint len) {
	return make_span_u8(ptr, len);
}

static span_u8
span_u8_slice_head(span_u8 s, uint n) {
	must(n <= s.len);
	return make_span_u8(s.ptr, n);
}

static span_u8
span_u8_slice_tail(span_u8 s, uint n) {
	must(n <= s.len);
	return make_span_u8(s.ptr + n, s.len - n);
}

static str
str_slice(str s, uint i, uint j) {
	must(i <= j);
	must(j <= s.len);
	return make_span_u8(s.ptr + i, j - i);
}

static str
str_slice_head(str s, uint n) {
	return span_u8_slice_head(s, n);
}

static str
str_slice_tail(str s, uint n) {
	return span_u8_slice_tail(s, n);
}

/*
Copies {n} bytes of memory from source to destination. Number of copied bytes
must be greater than zero.

Do not use for overlapping memory regions.

Params:
	{dst} - array pointer, must be not nil
	{src} - array pointer, must be not nil
	{n}   - must be not 0
*/
static void
unsafe_copy(u8* dst, u8* src, uint n) {
	for (uint i = 0; i < n; i += 1) {
		dst[i] = src[i];
	}
}

// Same as unsafe_copy, but bytes are copied in reverse order.
static void
unsafe_reverse_copy(u8* dst, u8* src, uint n) {
	uint i = 0;
	uint j = n - 1;
	while (i < n) {
		dst[j] = src[i];
		i += 1;
		j -= 1;
	}
}

/*/doc

Decode u32 integer from binary big endian encoding.
Given pointer must contain at least 4 bytes of memory.
*/
static u32
unsafe_get_u32be(u8* p) {
	u32 r = 0;
	
	// most significant byte 4
	r = p[0];

	// byte 3
	r <<= 8;
	r |= p[1];

	// byte 2
	r <<= 8;
	r |= p[2];

	// byte 1
	r <<= 8;
	r |= p[3];
	
	return r;
}

// Returns number of bytes copied.
static uint
copy(span_u8 dst, span_u8 src) {
    uint n = min_uint(dst.len, src.len);
    if (n == 0) {
        return 0;
    }

    must(dst.ptr != nil);
    must(src.ptr != nil);
	must(dst.ptr != src.ptr);

    unsafe_copy(dst.ptr, src.ptr, n);
	return n;
}

static void
clear(span_u8 s) {
	for (uint i = 0; i < s.len; i += 1) {
		s.ptr[i] = 0;
	}
}

static uint
unsafe_get_c_string_length(const u8* ptr) {
	uint i = 0;
	while (ptr[i] != 0) {
		i += 1;
	}
	return i;
}

static c_string
unsafe_make_c_string_from_ptr(u8* ptr) {
	uint len = unsafe_get_c_string_length(ptr);
	c_string s = {};
	s.ptr = ptr;
	s.len = len;
	return s;
}

/*/doc

Parameter {max_len} specifies max length (in bytes) of the resulting string.

If null byte is not found before reaching {max_len}, then resulting string length
will be recorded as {max_len}.

Parameter {ptr} must be not nil.
*/
static str
make_str_from_c_string_ptr(u8* ptr, uint max_len) {
	must(ptr != nil);
	str s = {};
	s.ptr = ptr;

	for (uint i = 0; i < max_len; i += 1) {
		if (ptr[i] == 0) {
			s.len = i;
			return s;
		}
	}

	s.len = max_len;
	return s;
}

/*
Copy the string {s} into buffer {buf} and place null byte at the end.
Returns {c_string} sliced from {buf}.

Buffer must be not nil and {buf.len} >= {s.len} + 1. 
*/
static c_string
unsafe_copy_as_c_string(span_u8 buf, str s) {
	must(buf.len >= s.len + 1);
	must(buf.ptr != nil);

	unsafe_copy(buf.ptr, s.ptr, s.len);
	buf.ptr[s.len] = 0;
	
	c_string cstr = {};
	cstr.ptr = buf.ptr;
	cstr.len = s.len;
	return cstr;
}

/*
Returns true if two strings contain the same data.
*/
static bool
str_equal(str a, str b) {
	if (a.len != b.len) {
		return false;
	}
	if (a.ptr == b.ptr) {
		return true;
	}

	for (uint i = 0; i < a.len; i += 1) {
		if (a.ptr[i] != b.ptr[i]) {
			return false;
		}
	}
	return true;
}

static bool
str_has_prefix(str s, str prefix) {
	if (s.len < prefix.len) {
		return false;
	}
	return str_equal(str_slice_head(s, prefix.len), prefix);
}

static u64
djb2_hash64(span_u8 s) {
    u64 h = 5381; // magic number
    for (uint i = 0; i < s.len; i += 1) {
        h = ((h << 5) + h) + cast(u64, s.ptr[i]);
    }
    return h;
}

typedef struct {
	str* ptr;
	uint len;
} span_str;

static span_str
make_span_str(void* ptr, uint len) {
    span_str s = {};
	if (len == 0) {
		return s;
	}

    s.ptr = ptr;
    s.len = len;
    return s;
}

typedef struct {
	u64* ptr;
	uint len;
} span_u64;

static span_u64
make_span_u64(u64* ptr, uint len) {
    span_u64 s = {};
	if (len == 0) {
		return s;
	}

    s.ptr = ptr;
    s.len = len;
    return s;
}

typedef struct {
	uint* ptr;
	uint  len;
} span_uint;

static span_uint
make_span_uint(uint* ptr, uint len) {
    span_uint s = {};
	if (len == 0) {
		return s;
	}

    s.ptr = ptr;
    s.len = len;
    return s;
}

typedef struct {
	s64* ptr;
	uint len;
} span_s64;

static span_s64
make_span_s64(s64* ptr, uint len) {
    span_s64 s = {};
	if (len == 0) {
		return s;
	}

    s.ptr = ptr;
    s.len = len;
    return s;
}

static span_s64
span_s64_slice_head(span_s64 s, uint n) {
	must(n <= s.len);
	return make_span_s64(s.ptr, n);
}

static span_s64
span_s64_slice_tail(span_s64 s, uint n) {
	must(n <= s.len);
	return make_span_s64(s.ptr + n, s.len - n);
}

/*
Align integer {n} to an integer {a}.

Alignment {a} must be a power of 2.
*/
static uint
align_uint(uint n, uint a) {
	uint mask = a - 1;
	return (n + mask) & ~mask;
}

static uint
get_align_offset(void* ptr, uint a) {
	uint p = cast(uint, ptr);
	return align_uint(p, a) - p;
}

static u64
rotate_left_u64(u64 x, uint k) {
	const uint mask = 64 - 1;
	k = k & mask;
	if (k == 0) {
		return x;
	}

    return (x << k) | (x >> (64 - k));
}

/*/doc

Represents time duration of something with nanosecond precision.
Both positive and negative durations are possible.
*/
typedef struct {
	// Number of seconds.
	s64 sec;

	// Number of nanoseconds.
	s64 nsec;
} TimeDur;

/*/doc

Subtract two time durations from one another.
Returns equivalent of {a - b} operation.
*/
static TimeDur
time_dur_sub(TimeDur a, TimeDur b) {
	TimeDur sub;
	sub.sec = a.sec - b.sec;
	sub.nsec = a.nsec - b.nsec;
  	if (sub.nsec < 0) {
		sub.sec -= 1;
		sub.nsec += 1000000000; // number of nanoseconds in one second
	}
	return sub;
}

static u8
fmt_dec_digit(u8 x) {
	return x + cast(u8, '0');
}

static u8
fmt_hex_digit(u8 x) {
	if (x <= 9) {
		return fmt_dec_digit(x);
	}
	return x - 10 + cast(u8, 'A');
}

static void
unsafe_fmt_hex_byte(span_u8 buf, u8 x) {
	buf.ptr[0] = fmt_hex_digit(x >> 4);
	buf.ptr[1] = fmt_hex_digit(x & 0xF);
}

static void
unsafe_fmt_bin_byte(span_u8 buf, u8 x) {
	uint i = 8;
	while (i != 0) {
		i -= 1;
		buf.ptr[i] = fmt_dec_digit(x & 1);
		x >>= 1;
	}
}

const uint max_u64_hex_length = 16;

/*/doc

Formats a given u64 integer as a hexadecimal number of
fixed width (=16) format, prefixing significant digits with
zeroes if necessary. Buffer must be at least 16 bytes
long.
*/
static void
unsafe_fmt_hex_prefix_zeroes_u64(span_u8 buf, u64 x) {
	uint i = max_u64_hex_length;
	// digits are written from least to most significant bit
	while (i != 0) {
		i -= 1;
		const u8 n = cast(u8, x & 0xF);
		buf.ptr[i] = fmt_hex_digit(n);
		x = x >> 4;
	}
}

/*/doc

Returns number of bytes actually written. In practice it will always
be either 0 or 16.
*/
static uint
fmt_hex_prefix_zeroes_u64(span_u8 buf, u64 x) {
	if (buf.len < max_u64_hex_length) {
		return 0;
	}
	must(buf.ptr != nil);

	unsafe_fmt_hex_prefix_zeroes_u64(buf, x);
	return max_u64_hex_length;
}

const uint max_u32_hex_length = 8;

/*/doc

Formats a given u32 integer as a hexadecimal number of
fixed width (=8) format, prefixing significant digits with
zeroes if necessary. Buffer must be at least 16 bytes
long.
*/
static void
unsafe_fmt_hex_prefix_zeroes_u32(span_u8 buf, u32 x) {
	uint i = max_u32_hex_length;
	// digits are written from least to most significant bit
	while (i != 0) {
		i -= 1;
		const u8 n = cast(u8, x & 0xF);
		buf.ptr[i] = fmt_hex_digit(n);
		x = x >> 4;
	}
}

/*/doc

Returns number of bytes actually written. In practice it will always
be either 0 or 16.
*/
static uint
fmt_hex_prefix_zeroes_u32(span_u8 buf, u32 x) {
	if (buf.len < max_u32_hex_length) {
		return 0;
	}
	must(buf.ptr != nil);

	unsafe_fmt_hex_prefix_zeroes_u32(buf, x);
	return max_u32_hex_length;
}

/*/doc

Puts decimal digits of a number in reverse order inside an array of bytes.
Returns number of digits written.

Array must have enough space to hold produced bytes.

To clarify: last (least significant) digit is placed at {ptr[0]}.

This operation is safe if buffer is at least {max_u64_dec_length} (=20) bytes long.
*/
static uint
unsafe_fmt_reverse_dec_u64(u8* ptr, u64 x) {
	uint i = 0;
	while (true) {
		const u8 n = cast(u8, x % 10);
		x /= 10;
		ptr[i] = fmt_dec_digit(n);
		i += 1;
		if (x == 0) {
			return i;
		}
	}
}

/*/doc

Maximum number of bytes (digits) needed to format any u64 integer as decimal number.
*/
static const uint
max_u64_dec_length = 20;

/*
Same as fmt_dec_u64, but has no buffer boundary checks.
*/
static uint
unsafe_fmt_dec_u64(span_u8 buf, u64 x) {
	u8 digits[max_u64_dec_length];
	const uint n = unsafe_fmt_reverse_dec_u64(digits, x);
	unsafe_reverse_copy(buf.ptr, digits, n);
	return n;
}

static void
unsafe_fmt_reverse_dec_fixed_width_u64(u8* ptr, u64 x, uint w) {
	uint i = 0;
	while (i < w) {
		const u8 n = cast(u8, x % 10);
		x /= 10;
		ptr[i] = fmt_dec_digit(n);
		i += 1;
	}
}

/*/doc

Formats integer as decimal number with fixed number of digits.
If number has not enough digits to fill the whole width than
zeroes are prefixed. If number has more digits then specified width
than only least significant ones will be written.
*/
static void
unsafe_fmt_dec_fixed_width_u64(span_u8 buf, u64 x, uint w) {
	u8 digits[max_u64_dec_length];
	unsafe_fmt_reverse_dec_fixed_width_u64(digits, x, w);
	unsafe_reverse_copy(buf.ptr, digits, w);
}

/*/doc

Maximum number of bytes (digits) needed to format any u64 integer as decimal number.
*/
static const uint
max_s64_dec_length = max_u64_dec_length + 1;

static uint
unsafe_fmt_dec_s64(span_u8 buf, s64 x) {
	if (x >= 0) {
		return unsafe_fmt_dec_u64(buf, cast(u64, x));
	}

	buf.ptr[0] = '-';
	uint n = unsafe_fmt_dec_u64(span_u8_slice_tail(buf, 1), cast(u64, -x));
	return n + 1;
}

/*
Puts integer x into a given buffer as text in decimal format.

Returns number of bytes written. Returns 0 if buffer does not
have enough bytes for the operation.
*/
static uint
fmt_dec_u64(span_u8 buf, u64 x) {
	if (buf.len == 0) {
		return 0;
	}
	must(buf.ptr != nil);
	
	if (buf.len >= max_u64_dec_length) {
		return unsafe_fmt_dec_u64(buf, x);
	}
	
	u8 digits[max_u64_dec_length];
	const uint len = unsafe_fmt_reverse_dec_u64(digits, x);
	if (len > buf.len) {
		// Not enough space in buffer to represent integer.
		return 0;
	}

	unsafe_reverse_copy(buf.ptr, digits, len);
	return len;
}

static uint
fmt_dec_s64(span_u8 buf, s64 x) {
	if (x >= 0) {
		return fmt_dec_u64(buf, cast(u64, x));
	}

	if (buf.len < 2) {
		return 0;
	}

	const uint n = fmt_dec_u64(span_u8_slice_tail(buf, 1), cast(u64, -x));
	if (n == 0) {
		return 0;
	}
	buf.ptr[0] = '-';
	return n + 1;
}

static const uint
max_time_dur_micro_length = max_s64_dec_length + 1 + 6;

/*/doc

Formats time duration like this: 17.053793 (with microsecond precision).
Returns number of bytes actually written.
*/
static uint
unsafe_fmt_time_dur_micro(span_u8 buf, TimeDur t) {
	// TODO: maybe we should support negative durations as well
	uint n = unsafe_fmt_dec_u64(buf, cast(u64, t.sec));
	
	buf.ptr[n] = '.';
	n += 1;
	
	span_u8 tail = span_u8_slice_tail(buf, n);
	unsafe_fmt_dec_fixed_width_u64(tail, cast(u64, t.nsec) / 1000, 6);
	return n + 6;
}

typedef struct {
	// Pointer to buffer memory.
	// Always not nil for initialized struct.
	u8* ptr;

	// Content length. How many bytes is already written to buffer.
	// Length cannot be greater than capacity.
	uint len;

	// Buffer storage capacity. Max bytes it can hold.
	uint cap;
} FormatBuffer;

static void
init_fmt_buffer(FormatBuffer* buf, span_u8 s) {
	must(s.len != 0);
	must(s.ptr != nil);

	buf->ptr = s.ptr;
	buf->len = 0;
	buf->cap = s.len;
}

static void
fmt_buffer_reset(FormatBuffer* buf) {
	buf->len = 0;
}

static str
fmt_buffer_head(FormatBuffer* buf) {
	return make_str(buf->ptr, buf->len);
}

static span_u8
fmt_buffer_tail(FormatBuffer* buf) {
	return make_str(buf->ptr + buf->len, buf->cap - buf->len);
}

static void
unsafe_fmt_buffer_put_str(FormatBuffer* buf, str s) {
	if (s.len == 0) {
		return;
	}

	span_u8 tail = fmt_buffer_tail(buf);
	must(s.len <= tail.len);

	unsafe_copy(tail.ptr, s.ptr, s.len);
	buf->len += s.len;
}

static void
unsafe_fmt_buffer_put_byte(FormatBuffer* buf, u8 b) {
	span_u8 tail = fmt_buffer_tail(buf);
	must(tail.len >= 1);

	tail.ptr[0] = b;
	buf->len += 1;
}

static void
unsafe_fmt_buffer_put_byte_repeat(FormatBuffer* buf, u8 b, uint n) {
	span_u8 tail = fmt_buffer_tail(buf);
	must(tail.len >= n);

	for (uint i = 0; i < n; i += 1) {
		tail.ptr[i] = b;
	}
	buf->len += n;
}

static void
unsafe_fmt_buffer_put_hex_byte(FormatBuffer* buf, u8 x) {
	span_u8 tail = fmt_buffer_tail(buf);
	must(tail.len >= 2);

	unsafe_fmt_hex_byte(tail, x);
	buf->len += 2;
}

static void
unsafe_fmt_buffer_put_bin_byte(FormatBuffer* buf, u8 x) {
	span_u8 tail = fmt_buffer_tail(buf);
	must(tail.len >= 8);

	unsafe_fmt_bin_byte(tail, x);
	buf->len += 8;
}

static void
unsafe_fmt_buffer_put_newline(FormatBuffer* buf) {
	unsafe_fmt_buffer_put_byte(buf, '\n');
}

static void
unsafe_fmt_buffer_put_space(FormatBuffer* buf) {
	unsafe_fmt_buffer_put_byte(buf, ' ');
}

static void
unsafe_fmt_buffer_put_space_repeat(FormatBuffer* buf, uint n) {
	unsafe_fmt_buffer_put_byte_repeat(buf, ' ', n);
}

static void
unsafe_fmt_buffer_put_hex_span_u8(FormatBuffer* buf, span_u8 s) {
	if (s.len == 0) {
		return;
	}

	unsafe_fmt_buffer_put_hex_byte(buf, s.ptr[0]);
	for (uint i = 1; i < s.len; i += 1) {
		unsafe_fmt_buffer_put_space(buf);
		unsafe_fmt_buffer_put_hex_byte(buf, s.ptr[i]);
	}
}

static void
unsafe_fmt_buffer_put_dec_u64(FormatBuffer* buf, u64 x) {
	const uint n = unsafe_fmt_dec_u64(fmt_buffer_tail(buf), x);
	buf->len += n;
}

static void
unsafe_fmt_buffer_put_dec_span_uint(FormatBuffer* buf, span_uint s) {
	if (s.len == 0) {
		return;
	}

	unsafe_fmt_buffer_put_dec_u64(buf, s.ptr[0]);
	for (uint i = 1; i < s.len; i += 1) {
		unsafe_fmt_buffer_put_space(buf);
		unsafe_fmt_buffer_put_dec_u64(buf, s.ptr[i]);
	}
}

static void
unsafe_fmt_buffer_put_hex_prefix_zeroes_u64(FormatBuffer* buf, u64 x) {
	unsafe_fmt_hex_prefix_zeroes_u64(fmt_buffer_tail(buf), x);
	buf->len += max_u64_hex_length;
}

static void
unsafe_fmt_buffer_put_hex_prefix_zeroes_u32(FormatBuffer* buf, u32 x) {
	unsafe_fmt_hex_prefix_zeroes_u32(fmt_buffer_tail(buf), x);
	buf->len += max_u32_hex_length;
}

/*/doc

Error codes are mostly used for low-level APIs or as a part of more
complex error types.

Different kind of errors can have the same codes.
Zero value always means no error.
*/
typedef uint ErrorCode;

typedef struct {
	span_u8 span;

	// Keeps track of block origin and generation.
	// Used for runtime checks and debug purposes.
	u64 id;
} MemBlock;

static void
clear_mem_block(MemBlock* block) {
	block->span.ptr = nil;
	block->span.len = 0;
	block->id = 0;
}

#define MEM_BUMP_ALLOCATOR 1

typedef struct {
	void* ptr;
	uint  kind;
} MemAllocator;

typedef struct {
	// Base pointer of the allocator.
	// Remains unchanged during its lifetime.
	u8* base;

	// Current allocator position (offset) relative to base.
	// Equals number of bytes already in use.
	uint pos;

	// Number of bytes in the underlying block of memory.
	// Remains unchanged during its lifetime.
	uint limit;
} MemBumpAllocator;

static void
init_mem_bump_allocator(MemBumpAllocator* al, span_u8 buf) {
	must(buf.len != 0);
	must(buf.ptr != nil);

	al->base = buf.ptr;
	al->pos = get_align_offset(buf.ptr, 8);
	al->limit = buf.len;
}

static MemAllocator
imake_mem_bump_allocator(MemBumpAllocator* al) {
	MemAllocator iface = {};
	iface.ptr = al;
	iface.kind = MEM_BUMP_ALLOCATOR;
	return iface;
}

#define ERROR_NO_MEMORY 1

static ErrorCode
mem_bump_alloc(MemBumpAllocator* al, MemBlock* block) {
	uint n = block->span.len;
	n = align_uint(n, 8);
	uint left = al->limit - al->pos;
	if (left < n) {
		return ERROR_NO_MEMORY;
	}

	block->span.len = n;
	block->span.ptr = al->base + al->pos;
	block->id = al->pos;

	al->pos += n;
	return 0;
}

/*
Request memory block from allocator.

Requested number of bytes (x) is specified by setting the field:

	block.span.len = x;

Allocator returns 0 on successfully allocating a block. In this case
allocator sets all fields in {block} argument.

	block.span.ptr - points to start (first byte) of the block
	block.span.len - actual number of bytes in allocated block, which is always not less than requested.
	block.id       - pseudo-unique identifier of the block

Number of requested bytes must be not zero.
*/
static ErrorCode
mem_alloc(MemAllocator al, MemBlock* block) {
	must(al.ptr != nil);
	must(block->span.len != 0);

	switch (al.kind) {
	case 0:
		panic_trap();
	case MEM_BUMP_ALLOCATOR:
		return mem_bump_alloc(al.ptr, block);
	default:
		panic_trap();
	}
}

static void
mem_free(MemAllocator al, MemBlock block) {
	if (block.span.len == 0) {
		// Freeing empty block must be nop
		return;
	}
	must(al.ptr != nil);
	must(block.span.ptr != nil);

	switch (al.kind) {
	case 0:
		panic_trap();
	case MEM_BUMP_ALLOCATOR:
		// Bump allocator does not free memory blocks
		return;
	default:
		panic_trap();
	}
}

#define PROC_MEM_BUFFER_SIZE 1 << 24

static u8 proc_mem_buffer[PROC_MEM_BUFFER_SIZE];
static MemBumpAllocator proc_mem_bump_allocator;

static void
init_proc_mem_bump_allocator() {
	init_mem_bump_allocator(&proc_mem_bump_allocator, make_span_u8(proc_mem_buffer, PROC_MEM_BUFFER_SIZE));
}

/*/doc

Global allocation mechanism which is meant for things that are created once
during process lifetime (and most likely live until process exits).
*/
static ErrorCode
proc_mem_alloc(MemBlock* block) {
	return mem_bump_alloc(&proc_mem_bump_allocator, block);
}

typedef struct {
	str name;
	str value;
} EnvEntry;


typedef struct {
	EnvEntry* ptr;
	uint len;
} span_env;

typedef struct {
	span_str args;
	span_env envs;
} OsProcInput;

OsProcInput os_proc_input;

static ErrorCode
init_os_proc_input(uint argc, u8** argv, u8** envp) {
	MemBlock block;
	block.span.len = argc * sizeof(str);
	ErrorCode code = proc_mem_alloc(&block);
	if (code != 0) {
		return code;
	}

	os_proc_input.args = make_span_str(block.span.ptr, argc);
	for (uint i = 0; i < argc; i += 1) {
		u8* ptr = argv[i];
		if (ptr == nil) {
			break;
		}
		os_proc_input.args.ptr[i] = unsafe_make_c_string_from_ptr(ptr);
	}

	return 0;
}
