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

typedef u64 uint;
typedef s64 sint;

typedef float      f32;
typedef double     f64;
typedef __float128 f128;

// This should only be used with pointer types.
#define nil 0

typedef struct {
    // Array pointer to raw bytes.
    u8* ptr;

    // Number of bytes available in {ptr}.
    uint num;
} s_u8, str;

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
