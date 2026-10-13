# Variables and constants

When it comes to identifiers that carry values various
languages usually categorise them something like:

- compile-time constants
- runtime constants (can be assigned only once)
- mutable variables (value can be reassigned multiple times)
- function parametes (some languages allow them to be mutated while others not)
- rebindable names (mostly in dynamically typed languages, Rust also does variation of this)

Ku follows somewhat similar approach, although prefers to separate them
based on usage intent and derive their properties from there.
Because of this it may be called "backwards" in logic. That is fine,
we are aware (and even proud) of it.

Ku value-carrying identifiers are separated into 3 disctinct categories:

## Constants

These are things that programmer would like to name once and refer to
them relying on that their value would never change. Thus we derive
that constants in Ku language are compile-time known constants and cannot
be changed during program execution by any means.

## Variables

These are actual variables and can be mutated freely at runtime as
their name suggest (values they hold can vary, change).

## Transfers

These are more akin to named containers for values. Values they hold
can be changed but only in restricted manner. Transfers are meant to
carry values from place to place without changing them. After value was
transfered and no longer in need, new value can be stored in container.
One may argue that difference between actual variables and transfers
are purely philosophical at this point and that may be so if we discard
different sets of allowed operations on them. But even if that is
the case, so what? We already stated that we derive separation of 
value-carriers based on usage intent, not their orthogonality of behaviour.

To summarize both variables and transfers are runtime variables in
traditional sense, but transfers allow less mutating operations and
require special constructs for them to happen.

Variables are meant to be used for implentation of algorithms and
functions where mutation is natural, needed, essential, leads to
clearer logic and less name clutter. Good usecases for variables are
demonstrated in example below, which implements a function for searching
maximum integer in a given span.

```ku
// returns (x, true) if span contains at least one element
// where x is >= all other elements
//
// returns (0, false) if span is empty and thus maximum cannot
// be found
fun get_max(s: []u32) -> (u32, bool) {
    if s.num == 0 {
        return 0, false;
    }

    // variable m holds current maximum during scan into span
    var m := s[0];

    // x is a transfer and cannot be changed in loop body
    //
    // here we declare intent clearly, that x should be managed
    // by loop automatically, we cannot accidentally mutate it
    for let x in s[1:] {
        if x > m {
            m = x; // m can be mutated because it is a variable
        }
    }
    return m, true;
}
```

Transfers should be used for "plumbing" values between expression
and calls. Consider example for typical logic of reading config file.

```
fun load_config(name: str) -> (*Config, error) {
    const dir_path = "env/configs/";

    // construct "let" creates new transfer(s)
    let path := strings.join(dir_path, name);

    // here it creates two transfers "file" and "err"
    let file, err := os.open(path);
    if err != nil {
        // ... logic for error handling
    }

    // transfers can be given a new value via "~" operator,
    // each separate entry in "let" construct must either be
    // a new transfer name or existing one with explicit "~"
    // prefix which specifies that we want to change value
    // stored in it
    //
    // transfer type is fixed, it cannot be changed even via
    // "let" + "~" construct (unlike in Rust for example)
    let data, ~err := io.read_all(file);
    if err != nil {
        // ... logic for error handling
    }

    // specifying "~" for each "err" after the first is annoying
    // and maybe rules will be more relaxed in the future specifically
    // for "error" and "errid" types if they come last in multiple
    // return values
    let config, ~err := parse_config(data);
    if err != nil {
        // ... logic for error handling
    }

    return config, nil;
}
```
