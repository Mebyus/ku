/*/doc

Removes space characters (' ' byte) from both start and end of the string.
Returned string is either:
1. empty in case original string contains only spaces or is empty
2. slice into original string
*/
static str
str_trim_space(str s) {
    uint i = 0;
    while (i < s.len && s.ptr[i] == ' ') {
        i += 1;
    }
    str tail = str_slice_tail(s, i);

    i = tail.len;
    while (i != 0) {
        i -= 1;
        if (tail.ptr[i] != ' ') {
            return str_slice_head(tail, i + 1);
        }
    }
    return empty_str;
}

/*/doc

Finds first line break and slices string head before it.
If there is no line break in the string then entire string is returned. 
*/
static str
str_slice_line(str s) {
    RetIndex r = str_index_byte(s, '\n');
    if (!r.ok) {
        return s;
    }
    return str_slice_head(s, r.index);
}

#define SMALL_STRING_PATTERN (1 << 10)

/*/doc

Argument {lps} must have the same length as string {p}.
Length must be at least 1.
*/
static void
unsafe_str_fill_longest_prefix_suffix(span_uint lps, str p) {
    // Length of the previous longest prefix suffix
    uint k = 0;

    lps.ptr[0] = 0;

    uint i = 1;
    while (i < p.len) {
        if (p.ptr[i] == p.ptr[k]) {
            k += 1;
            lps.ptr[i] = k;
            i += 1;
        } else {
            if (k != 0) {
                k = lps.ptr[k - 1];
            } else {
                lps.ptr[i] = 0;
                i += 1;
            }
        }
    }
}

/*/doc

Function uses {lps} argument as a buffer. Its length must be the same as {p} length.
Uses KMP (Knuth-Morris-Pratt) algorithm.
*/
static RetIndex
unsafe_str_index_kmp(str s, str p, span_uint lps) {
    unsafe_str_fill_longest_prefix_suffix(lps, p);
    
    uint i = 0; // index for {s}
    uint j = 0; // index for {p}
    
    RetIndex ret = {};
    while ((s.len - i) >= (p.len - j)) {
        if (s.ptr[i] == p.ptr[j]) {
            i += 1;
            j += 1;
        }

        if (j >= p.len) {
            // match found
            ret.index = i - j;
            ret.ok = true;
            return ret;
        }

        if (i < s.len && s.ptr[i] != p.ptr[j]) {
            if (j != 0) {
                j = lps.ptr[j - 1];
            } else {
                i += 1;
            }
        }
    }

    // no match found
    return ret;
}

static RetIndex
unsafe_str_index_small(str s, str p) {
    must(p.len <= SMALL_STRING_PATTERN);

    uint buf[SMALL_STRING_PATTERN];
    span_uint lps = make_span_uint(buf, p.len);
    return unsafe_str_index_kmp(s, p, lps);
}

/*/doc

Searches for a string pattern {p} inside another string {s}.
If match is found, returns index of the first match such that
s[i] == p[0] is true.

Argument {p} must not be an empty string.
*/
static RetIndex
str_index(str s, str p) {
    must(p.len != 0);

    RetIndex ret = {};
    if (p.len > s.len) {
        return ret;
    }

    if (p.len < SMALL_STRING_PATTERN) {
        return unsafe_str_index_small(s, p);
    }

    // TODO: code for patterns that are not small, probably need an allocator here
    return ret;
}

typedef struct {
    str  prefix;
    uint count;
    uint weight;
} RetPrefixRepeat;

/*/doc

Counts how many times {prefix} string repeats (with no overlaps)
at the beginning of string {s}.

If string {s} does not start with {prefix}, than returned count is 0.
*/
static uint
str_count_prefix_repeats(str s, str prefix) {
    if (s.len < prefix.len) {
        return 0;
    }

    uint count = 0;
    uint i = 0;
    while (i < s.len) {
        if (!str_has_prefix(str_slice_tail(s, i), prefix)) {
            return count;
        }

        i += prefix.len;
        count += 1;
    }
    return count;
}

static RetPrefixRepeat
str_find_optimal_prefix_repeats(str s) {
    RetPrefixRepeat ret = {};
    if (s.len == 0) {
        return ret;
    }

    // first we check repeats of one leading byte by a simple loop
    uint i = 1;
    while (i < s.len) {
        if (s.ptr[i] == s.ptr[0]) {
            i += 1;
        } else {
            break;
        }
    }

    // best total length of prefix repeats = prefix.len * count
    uint best_count = i;
    uint best_total_len = i;
    uint best_prefix_len = 1;
    uint best_weight = best_total_len - best_prefix_len;

    // current prefix length
    uint len = 2;
    while (len <= (s.len >> 1)) {
        uint count = 1 + str_count_prefix_repeats(str_slice_tail(s, len), str_slice_head(s, len));
        uint total_len = len * count;
        uint weight = total_len - len;

        if (count > 1 && weight > best_weight) {
            best_total_len = total_len;
            best_count = count;
            best_prefix_len = len;
            best_weight = weight;
        }

        len += 1;
    }

    ret.prefix = str_slice_head(s, best_prefix_len);
    ret.count = best_count;
    ret.weight = best_weight;
    return ret;
}
