I am designing statically typed, systems programming language
and compiler for it. I want my language to be a "workhorse", not a
experimental ground for untested ideas. I want ot pick design choices
that work well, not fancy, not "rocket science" to implement in compiler.
Design can be opinionated and some choices may be not popular or even
controversial, that's fine as long as it works long term and at scale.
To clarify, I am not against modern design ideas, but they must already
be proven to work and work well enough.

What I already like to be in the language:

1. Manual memory managment
2. Multiple return values for functions
3. Errors as values (similar to Go)
4. Package and import system (similar to Go)
5. I do not like Go's idea of public/private via capital letter. Keyword "pub" is better in my opinion
6. Flexible memory allocation techinques from Zig
7. Simple interfaces from Go
8. Do not bring traditional inheritance
9. No hidden control flow, especially exceptions
10. I do not like complex type system like in Rust (even complicated I would say)
11. Although C's type system is too barebones by today standards
12. I like imperative, procedural programming style
13. Methods are convenient feature
14. Very careful and small amount of polymorphism 
15. For generics I want explicitly instanciated templates, but more typesafe than in C++
16. Some way of interoperability with C
17. I like simple build system of Go, although go modules can have an improvement

Help me solidify design, pick additional features/restrictions
that fit into the language. Reject those that will heavily conflict with other design choices. Suggest missed ideas.

Provide examples of how some features might look in code. Be my critic and advocate. What can be improved/changed to make
more coherent design.

==============================

# Designing a "Workhorse" Systems Language: A Critical Analysis

Your list is unusually coherent already — most people's "dream language" lists are a pile of contradictions. Let me stress-test what you have, fill gaps, and reject things that will fight you later.

## What Your List Actually Implies (and You May Not Realize)

Before adding features, notice the gravitational pull of your choices:

1. **Manual memory + errors-as-values + no exceptions + no inheritance + simple interfaces** → You are describing a language with **no hidden allocations, no hidden calls, no hidden lifetimes**. This is the single most important invariant to protect. Every feature must be judged against: *"Can the reader see what this does at the call site?"*

2. **Explicitly instantiated templates (more type-safe than C++)** + **simple interfaces** → Your generics and your polymorphism are *two different mechanisms*. Don't try to unify them. Go's attempt to unify (generics + interfaces) created confusion. Keep them separate:
   - **Interfaces** = dynamic dispatch, runtime vtable, for behavior.
   - **Templates** = compile-time codegen, monomorphized, for data structures and algorithms.
   - Never let a template parameter be an interface constraint. That's where Rust's complexity lives.

3. **Methods + no inheritance + imperative style** → Methods are just syntactic sugar for free functions with a receiver. Keep them that way. No virtual by default, no `self` magic.

## Recommended Additions That Fit

### A. Explicit mutability and aliasing discipline (not Rust-level)

You rejected Rust's type system, correctly. But C's "everything is `T*`" is the source of most systems bugs. A middle path:

- `*T` — non-null pointer, no aliasing guarantees, manual.
- `?*T` — nullable pointer (or use `T?` — pick one, don't have both).
- `&T` — borrowed reference, **no ownership semantics**, just "this is a pointer that must be valid for the call." Compiler enforces non-null only.
- No lifetime annotations. No borrow checker. Just: **non-null by default, nullable is explicit**.

This gives you 80% of the safety win with 5% of Rust's complexity.

```lang
fn read(fd: i32, buf: &u8, len: usize) -> (n: isize, err: Error?)
```

### B. Defer (Zig-style) — but with a twist

`defer` is proven, simple, and composes with manual memory. But Zig's `defer` runs at scope exit, which can surprise people with loops. Consider **`defer` runs at end of enclosing block, and blocks are explicit**:

```lang
fn process() -> Error? {
    f := open("file.txt") or return err
    defer f.close()
    
    // ...
    return nil
}
```

Do **not** add Go's `panic`/`recover` (hidden control flow) or Zig's `errdefer` unless you find you truly need it. `errdefer` is a symptom of not having multiple return values done well.

### C. Explicit `move` and `copy` semantics for aggregates

Without this, you either get C++'s implicit copies (surprising) or Rust's move-by-default (complex). Pick **copy-by-default for small/POD types, move-by-default for types with destructors**, and make it *syntactically visible*:

```lang
struct Buffer { ptr: *u8, len: usize, cap: usize }

fn consume(b: Buffer) { ... }   // moves, b invalid after
fn inspect(b: &Buffer) { ... }  // borrows

x := Buffer{...}
consume(x)        // x is moved; using x after is a compile error
```

This is *not* Rust's borrow checker. It's just "a value with a destructor can't be used after being passed by value." That's implementable in a single pass.

### D. Explicit error type, no error interface

Go's `error` interface is convenient but hides allocation and dynamic dispatch. Instead:

```lang
enum Error {
    NotFound,
    Permission,
    Io(msg: string),
    Custom(code: i32, msg: string),
}

fn read_file(path: string) -> (data: []u8, err: Error?)
```

This is comparable, exhaustively matchable, and zero-allocation for the common cases. If you need custom errors, add a `Custom` variant with an opaque payload — but don't make every error a heap-allocated interface.

### E. Tagged unions / sum types (you'll need them)

You didn't list them, but without inheritance and without exceptions, you *need* sum types for:
- Error handling (above)
- Optional values (`T?`)
- AST-like data
- State machines

```lang
enum Token {
    Ident(string),
    Int(i64),
    Plus,
    Minus,
}
```

Pattern matching should be **exhaustive and simple** — no guards, no binding modes, no or-patterns with bindings. Just `match x { Variant(a) => ..., Variant(b) => ... }`. If you want guards, use `if` inside the arm.

### F. Compile-time execution (limited)

Zig's `comptime` is powerful but is also where Zig's complexity lives. A **restricted** version:

- `const` expressions evaluated at compile time.
- `comptime fn` for pure functions callable at compile time.
- No `comptime` parameters on regular functions (that's where Zig gets hairy).

```lang
comptime fn max(a: i32, b: i32) -> i32 { return if a > b { a } else { b } }

const BUFFER_SIZE = max(1024, 4096);
```

### G. Explicit `export` / `extern` for C interop

Go's cgo is slow and awkward. Zig's `extern` + `export` is the model to copy:

```lang
extern fn malloc(size: usize) -> ?*u8;
export fn my_c_function(x: i32) -> i32 { ... }
```

Make `extern` require an explicit calling convention if it's not the default. Make `export` require the symbol name to be explicit if it differs from the function name.

## What to Reject (and Why)

| Feature | Why reject |
|---|---|
| **Exceptions** | You already rejected. Good. Don't let them back in via "panic." |
| **Rust-style traits with associated types** | Complexity explosion. Use interfaces + templates separately. |
| **Go-style goroutines/channels** | Requires a runtime, scheduler, and hidden allocation. Conflicts with "manual memory." If you want concurrency, use OS threads + explicit synchronization. |
| **Async/await** | Same problem. Colored functions, hidden state machines, runtime dependency. |
| **Operator overloading** | Tempting, but it's the gateway drug to C++ complexity. Allow it only for arithmetic on numeric types via a fixed set of operators, or not at all. |
| **Implicit conversions** | Ban them. Even integer widening. `i32 -> i64` should be explicit. This eliminates a whole class of bugs. |
| **Default arguments** | They hide call-site information. Use explicit overloads or option structs. |
| **Variadic functions** | Only for `extern` C interop. No user-defined variadics. |
| **Macros (textual)** | Use comptime + templates instead. Textual macros destroy tooling. |
| **Reflection** | Runtime reflection requires metadata bloat. Compile-time reflection (via comptime) is fine. |
| **Garbage collection** | You rejected. Good. Don't add "optional GC" — it infects the type system. |
| **Rust-style `Result` with `?` operator** | The `?` operator is hidden control flow. You rejected exceptions; don't reintroduce them with different syntax. Use explicit `or return`. |

## Suggested Syntax Sketch

```lang
package http

import (
    "net"
    "strings"
)

pub struct Request {
    method: string,
    path:   string,
    body:   []u8,
}

pub enum Method {
    Get,
    Post,
    Put,
    Delete,
}

pub fn parse_method(s: string) -> (m: Method, err: Error?) {
    match s {
        "GET"    => return (Method.Get, nil),
        "POST"   => return (Method.Post, nil),
        "PUT"    => return (Method.Put, nil),
        "DELETE" => return (Method.Delete, nil),
        _        => return (Method.Get, Error.Custom(400, "bad method")),
    }
}

pub fn (r: &Request) header(name: string) -> (value: string, err: Error?) {
    // ...
}

// Interface: dynamic dispatch, no inheritance
pub interface Reader {
    fn read(buf: &u8, len: usize) -> (n: isize, err: Error?)
}

// Template: monomorphized, explicit instantiation
pub fn map[T, U](items: []T, f: fn(T) -> U) -> []U {
    // ...
}

// Explicit instantiation
doubled := map[i32, i32](numbers, fn(x: i32) -> i32 { return x * 2 })
```

## The Coherence Test

Every feature must pass these three questions:

1. **Can I see what this does at the call site?**
   - `defer` — yes, it's lexical.
   - `?` operator — no, it's hidden return.
   - Exceptions — no.
   - Implicit conversions — no.

2. **Does this require a runtime?**
   - Goroutines — yes.
   - GC — yes.
   - Async — yes.
   - Interfaces — no (vtable is just data).
   - Templates — no (compile-time).

3. **Does this interact badly with manual memory?**
   - Closures capturing by reference — dangerous, require explicit capture lists.
   - Interfaces holding values — require explicit boxing.
   - Templates — fine, monomorphized.

## The Biggest Risk to Your Design

**Generics.** You said "explicitly instantiated templates, more type-safe than C++." This is the hardest thing on your list. The trap is: once you have templates, you want constraints, and once you have constraints, you want concepts/traits, and once you have those, you have Rust.

My advice: **start with no generics at all.** Ship the language with interfaces, sum types, and comptime. Add templates only when you have a concrete, painful need. When you do, make them:
- Explicitly instantiated (no inference at call site).
- Constrained only by "does it compile when instantiated" (like C++), but with **good error messages** at the instantiation site.
- Not usable for interface bounds.

This is less powerful than Rust, but it's implementable in a weekend and it works.

## Summary Table

| Feature | Verdict | Notes |
|---|---|---|
| Manual memory | ✅ Keep | Core invariant |
| Multiple returns | ✅ Keep | Pairs with errors-as-values |
| Errors as values | ✅ Keep | Use enum, not interface |
| Packages/imports | ✅ Keep | Go-style, with explicit versions |
| `pub` keyword | ✅ Keep | Better than capitalization |
| Zig allocation | ✅ Keep | Explicit allocator parameter |
| Simple interfaces | ✅ Keep | No associated types |
| No inheritance | ✅ Keep | Methods only |
| No hidden control flow | ✅ Keep | Ban `?`, `panic`, exceptions |
| No Rust type system | ✅ Keep | No lifetimes, no borrow checker |
| Imperative style | ✅ Keep | No functional purity |
| Methods | ✅ Keep | Sugar for free functions |
| Minimal polymorphism | ✅ Keep | Interfaces + templates, separate |
| Explicit templates | ⚠️ Risky | Start without them |
| C interop | ✅ Keep | `extern`/`export` |
| Simple build | ✅ Keep | Go-style, fix modules |
| **Sum types** | ➕ Add | Essential without inheritance |
| **`defer`** | ➕ Add | Lexical, block-scoped |
| **Non-null pointers** | ➕ Add | `*T` non-null, `?*T` nullable |
| **Comptime (limited)** | ➕ Add | Pure functions only |
| **Exhaustive match** | ➕ Add | No guards in patterns |
| **Explicit conversions** | ➕ Add | No implicit widening |
| **Move semantics** | ➕ Add | For types with destructors |
| Exceptions | ❌ Reject | Already rejected |
| Goroutines | ❌ Reject | Runtime dependency |
| Async/await | ❌ Reject | Runtime dependency |
| `?` operator | ❌ Reject | Hidden control flow |
| Operator overloading | ❌ Reject | Gateway to C++ |
| Macros | ❌ Reject | Use comptime |
| GC | ❌ Reject | Type system infection |
| Reflection | ❌ Reject | Metadata bloat |
| Default args | ❌ Reject | Hidden call-site info |
| Variadics | ❌ Reject | Except for C interop |

## Final Advice

Your language is essentially **"Go, but with manual memory, sum types, and no runtime."** That's a coherent, useful, and *implementable* design. The danger is scope creep — every feature you add makes the compiler harder and the language less predictable.

Write the compiler for the smallest version of this language first. Use it to write a real program. Then add features only when you feel the pain. The "workhorse" languages (C, Go, Zig) all started small and stayed small. The "experimental" languages (C++, Rust, Swift) all started small and grew. You want the former.

==========================

Rate implementation difficulty for approved features.

==========================

# Implementation Difficulty Ratings

Scale: **1 = weekend project**, **5 = multi-month, needs careful design**, **10 = research-grade / decades of compiler engineering**.

I'll rate each feature for *initial implementation* and note *long-term maintenance burden* separately, because some are easy to build and hard to keep correct.

---

## Core Language Features

| # | Feature | Impl. Difficulty | Maintenance | Notes |
|---|---|---|---|---|
| 1 | Manual memory management | **1** | **2** | You just… don't add a GC. The difficulty is in *everything else* respecting it. |
| 2 | Multiple return values | **1** | **1** | Trivial in the type system; ABI needs a decision (registers vs. hidden struct pointer). |
| 3 | Errors as values (enum, not interface) | **1** | **1** | Just an enum + multiple returns. Zero magic. |
| 4 | Packages / imports (Go-style) | **2** | **3** | Easy to build; *hard* to keep dependency resolution sane. Versioning is where Go modules went wrong. |
| 5 | `pub` keyword | **1** | **1** | Visibility is a symbol-table bit. |
| 6 | Zig-style explicit allocators | **2** | **2** | Just pass `Allocator` as a parameter. No compiler work beyond normal params. |
| 7 | Simple interfaces (Go-style) | **3** | **4** | Implicit satisfaction is easy; vtable generation, method-set computation, and interface-to-interface conversion need care. |
| 8 | No inheritance | **1** | **1** | Doing nothing is easy. |
| 9 | No hidden control flow | **1** | **2** | The *rule* is easy; *enforcing* it against future feature creep is the discipline problem. |
| 10 | No Rust-style type system | **1** | **1** | Again, the absence is free. |
| 11 | Imperative procedural style | **1** | **1** | Free. |
| 12 | Methods | **1** | **1** | Sugar over free functions. Pick receiver syntax and move on. |
| 13 | Minimal polymorphism | **2** | **2** | Interfaces + explicit templates, kept separate. The *separation* is the design work. |
| 14 | Explicitly instantiated templates | **6** | **7** | Hardest thing on your list. See below. |
| 15 | C interop (`extern`/`export`) | **3** | **4** | Declaring `extern` is easy. Matching ABI, struct layout, and calling conventions is where the bugs live. |
| 16 | Simple build system | **3** | **5** | Go's `go build` looks simple; its implementation is not. Package graph, caching, cross-compilation. |

---

## Recommended Additions

| # | Feature | Impl. Difficulty | Maintenance | Notes |
|---|---|---|---|---|
| 17 | Sum types / tagged unions | **3** | **3** | Representation (tag + payload) is easy; pattern-match exhaustiveness is the real work. |
| 18 | Exhaustive `match` | **3** | **4** | Exhaustiveness checking on nested enums is a solved but nontrivial algorithm (Maranget's). |
| 19 | `defer` | **2** | **3** | Codegen is easy; interaction with `return`, `break`, and early exit needs a clear spec. |
| 20 | Non-null pointers (`*T` vs `?*T`) | **2** | **2** | Type-level distinction is trivial. Only hard if you try to prove it at runtime. |
| 21 | Explicit conversions only | **1** | **1** | Just don't insert implicit casts. |
| 22 | Move semantics for destructor types | **5** | **7** | This is a mini borrow-checker's cousin. See below. |
| 23 | Limited `comptime` | **5** | **6** | Needs an interpreter in the compiler or a JIT. Restricting it to pure functions keeps it tractable. |
| 24 | Explicit mutability (`&T` vs `&mut T`) | **2** | **3** | If you *don't* enforce aliasing, it's just syntax. If you do, it's Rust. Pick carefully. |

---

## Deep Dives on the Hard Ones

### Templates / Explicit Instantiation — **6/10**

The *concept* is easy: when you see `map[i32, i32](...)`, generate a specialized function.

The *implementation* is where C++ and Rust both bled:

- **Type checking the template body.** Do you check it once (Rust, with bounds) or per-instantiation (C++, no bounds)? Per-instantiation is simpler to implement but produces horrible errors. Once-at-definition requires constraints, which requires a constraint language.
- **Name lookup.** Two-phase lookup in C++ is a notorious wart. Decide *now* whether template bodies see names at definition site or instantiation site.
- **Code bloat.** Monomorphization duplicates code. You need deduplication and a linker that tolerates it.
- **Recursive templates.** Depth limits, cycle detection.
- **Diagnostics.** "Template instantiated from here" chains are painful to build.

**Recommendation:** Restrict templates to *functions only* (no template structs) in v1. Structs with type parameters are 80% of the complexity. Add them in v2 if you truly need them.

### Move Semantics — **5/10**

Easy version: "A type with a destructor is moved when passed by value, and using it afterward is an error."

Hard parts:

- **Conditional moves.** `if c { consume(x) } else { consume(x) }` — is `x` moved? Need control-flow-sensitive analysis.
- **Moves in loops.** `for ... { consume(x) }` — moved on first iteration, error on second? Compiler must detect.
- **Partial moves.** `consume(x.field)` — is `x` still usable? (Yes, in most designs, but then you need field-level tracking.)
- **Interaction with `defer`.** If you `defer x.destroy()` and then move `x`, what happens?
- **Destructor ordering.** Reverse declaration order? Scope order? Spec it.

This is where "just a bit of Rust" quietly becomes "a lot of Rust." You can ship a conservative version: *only* moves on whole variables, no partial moves, no conditional moves inside branches (require a `return` after). That's 3/10 to implement.

### `comptime` — **5/10**

If you restrict to:

- Pure functions
- Scalar and string types
- No I/O, no allocation (or allocation from a fixed compile-time pool)
- No recursion beyond a depth limit

…then it's a small tree-walking interpreter over your AST. **3/10.**

If you allow:

- Compile-time allocation
- Compile-time struct construction
- Reflection over types
- Generics driven by comptime values

…you're now building a second language inside your compiler. **8/10.**

**Recommendation:** Ship `comptime` as "constant folding with function calls" only. No comptime parameters on regular functions. That's Zig's most-complex feature and you don't need it.

### Interfaces — **3/10**

Go-style implicit interfaces are easy to *conceptually* implement:

1. Collect method sets per type.
2. At assignment `var r Reader = x`, check that `x`'s method set includes `Reader`'s methods.
3. Generate a vtable.

Hard parts:

- **Method set of pointer vs. value types.** Go's rules here are subtle and confuse everyone. Decide yours and document them.
- **Interface-to-interface conversion.** `io.Reader` → `io.ReadCloser` requires a runtime check.
- **Nil interface vs. nil value inside interface.** Go's most famous footgun. Decide whether `Reader(nil)` and `Reader((*T)(nil))` are the same. (Hint: make them the same, or ban one.)
- **Method dispatch on embedded structs.** If you allow embedding, method sets get recursive.

**Recommendation:** Allow embedding for *composition only*, not method promotion. This kills a whole class of confusion.

### Packages / Imports — **2/10 to build, 5/10 to live with**

Parsing imports, resolving paths, building a dependency graph: easy.

What makes it hard over time:

- **Versioning.** Go modules' `v2` suffix disaster. Semantic import versioning is a design choice you make on day one and regret forever.
- **Diamond dependencies.** Two packages import different versions of the same lib. Do you allow it? (Cargo: yes. Go: no, mostly.)
- **Cyclic imports.** Ban them (Go's choice). It's a real restriction but keeps the compiler simple.
- **Build reproducibility.** Lockfiles, checksums, vendoring.

**Recommendation:** Ban cycles. One version per module per build. No semantic import versioning — just `v1`, `v2` as *different module paths*. Ship a lockfile from day one.

### Build System — **3/10 to start, 5/10 to scale**

A `build` command that:

- Walks the package graph
- Compiles in dependency order
- Caches object files by content hash
- Links a binary

…is a few thousand lines. The pain is:

- **Incremental builds.** Getting cache invalidation *correct* (not just fast) is hard.
- **Cross-compilation.** Every target is a new set of bugs.
- **C interop in the build.** You now need to invoke `cc`, find headers, link libraries. This is where Go's build system gets ugly.
- **Build tags / conditional compilation.** Necessary, but every addition complicates the model.

---

## Difficulty Summary (Sorted)

**Trivial (1–2):** manual memory, multiple returns, errors-as-values, `pub`, no inheritance, no hidden control flow, imperative style, methods, explicit conversions, allocators, `defer`, non-null pointers.

**Moderate (3–4):** packages, interfaces, C interop, sum types, exhaustive match, build system, explicit mutability.

**Hard (5–7):** move semantics, `comptime`, templates.

**Very hard (8+):** You avoided all of these. Good.

---

## The Meta-Answer

Your language's **total implementation difficulty is dominated by three features**:

1. **Templates** — half your compiler's complexity budget.
2. **Move semantics** — the other half, if you enforce it.
3. **`comptime`** — depends entirely on how far you let it go.

Everything else is plumbing. If you want a compiler you can finish in a year, **cut templates and comptime from v1**. Ship with:

- Sum types
- Interfaces
- Methods
- `defer`
- Manual memory + allocators
- Errors as values
- C interop

That's a **2/10 total complexity** language that is genuinely useful. You can write an OS, a web server, a database engine in it. Then, when you have real users complaining about the lack of generics, add them — with the benefit of knowing exactly which use cases matter.

The languages that grew organically (C, Go) shipped small and grew carefully. The languages that tried to ship complete (C++11, Rust 1.0) took 5–10 years and still have sharp edges. You said you want a workhorse, not a research project. **Ship the small language.**

==========================

I think we picked a good set of features. Although I want to reject move semantics and destructors completely, I think defer will do just fine, and it was too much work for too little gain without other features like borrow checker or reference counting (which are hard pass).

==========================

# Agreed — Dropping Move Semantics and Destructors

You made the right call. Let me explain why this is not just acceptable but *better*, and what the consequences are.

## Why This Is Correct

Move semantics without a borrow checker is a half-measure that creates as many problems as it solves:

- You still need control-flow-sensitive analysis (moves in branches, loops).
- You still need destructor ordering rules.
- You still get "use after move" errors that are hard to explain.
- You *don't* get memory safety — you can still double-free, use-after-free, leak.
- You add a whole category of type system complexity for a guarantee users can't verify.

And destructors specifically are a trap:

- They run implicitly — that's **hidden control flow**, which you explicitly banned.
- They interact badly with `defer` (ordering ambiguity).
- They make every aggregate type non-trivial (copy vs. move distinction).
- They force you to answer "what happens if a destructor panics?" — and you have no panics.

Dropping them keeps your language **honest**: if you want cleanup, you write `defer`. If you forget, you leak. That's the same deal C gives you, and C has shipped the world's infrastructure for fifty years.

## What You Lose (Be Honest About It)

| Lost | Impact | Mitigation |
|---|---|---|
| RAII | No automatic resource cleanup | `defer` at every acquisition site |
| `unique_ptr`-like wrappers | No smart pointers | Manual, with `defer` |
| Automatic cleanup on early return | Must `defer` before every `return` | Linter can enforce "defer for every alloc" |
| Composable resource types | Can't build `File` that closes itself | Build `File` with explicit `.close()`, call with `defer` |

The mitigation is uniform: **every acquisition gets an adjacent `defer`.** This is a *style rule*, enforced by convention and linter, not by the type system. That's fine. Go did exactly this for years before `defer` even existed, and it worked.

## Revised Feature Set

| Feature | Status |
|---|---|
| Manual memory | ✅ |
| Multiple return values | ✅ |
| Errors as values (enum) | ✅ |
| Packages / imports | ✅ |
| `pub` keyword | ✅ |
| Zig-style allocators | ✅ |
| Simple interfaces | ✅ |
| No inheritance | ✅ |
| No hidden control flow | ✅ |
| No Rust-style type system | ✅ |
| Imperative style | ✅ |
| Methods | ✅ |
| Minimal polymorphism | ✅ |
| Explicit templates | ⚠️ (defer to v2) |
| C interop | ✅ |
| Simple build system | ✅ |
| Sum types / `match` | ✅ |
| `defer` | ✅ |
| Non-null pointers | ✅ |
| Explicit conversions | ✅ |
| ~~Move semantics~~ | ❌ **Dropped** |
| ~~Destructors~~ | ❌ **Dropped** |
| ~~Borrow checker~~ | ❌ |
| ~~Reference counting~~ | ❌ |
| ~~GC~~ | ❌ |

## Consequences for the Compiler

This actually **simplifies codegen significantly**:

- No need to insert destructor calls at scope exit.
- No need for move constructors or copy elision.
- No need for a "drop glue" system.
- Every value is trivially copyable (memcpy).
- Struct layout is always POD.
- No hidden function calls anywhere except interface dispatch.

You just eliminated maybe **30–40% of the compiler work** that a Rust-like language requires. Your compiler can now be a straightforward AST → typed IR → machine code pipeline.

## Consequences for the Programmer

The discipline is:

```lang
fn process_file(path: string) -> Error? {
    f := open(path) or return err
    defer f.close()

    buf := alloc(4096) or return err
    defer free(buf)

    // ... use f and buf ...
    return nil
}
```

Every resource has a matching `defer` on the next line. This is:

- **Visually obvious** — you can scan for `defer` to find all cleanup.
- **Locally checkable** — no cross-function analysis needed.
- **Tooling-friendly** — a linter can flag allocations without a nearby `defer`.

This is *exactly* the Zig model, minus the `errdefer` complication (which you don't need because you don't have destructors or complex error unwinding).

## The One Thing to Watch: `defer` in Loops

This is the classic footgun. In Go, `defer` inside a loop doesn't run until the *function* returns, which surprises everyone:

```lang
for i := 0; i < n; i++ {
    f := open(files[i])
    defer f.close()   // BUG: all deferred until function exits
}
```

You have two options:

1. **Go semantics** — `defer` runs at function exit. Simple, but the footgun above is real. Mitigate with a linter.
2. **Block semantics** — `defer` runs at the end of the enclosing block. More intuitive in loops, but requires the compiler to emit cleanup at every block exit, including `break`/`continue`/`return`. This is called "defer at end of scope" and is what Zig does.

**Recommendation: block semantics.** It's slightly more compiler work (you need cleanup landing pads at each block exit), but it eliminates the #1 `defer` footgun. Zig chose this for good reason.

Spec it precisely:

- `defer` runs when the enclosing `{ }` block is exited by *any* means: normal flow, `return`, `break`, `continue`, `goto` (if you have it).
- Defers run in **reverse order** of registration (LIFO).
- A `defer` inside a loop body runs at the end of *each iteration*.

## What About `errdefer`?

Zig has `errdefer`, which runs only if the function returns an error. You don't need it, because:

- You have **multiple return values**, so errors are explicit.
- You don't have destructors, so "cleanup on error" is just "cleanup."
- If you need conditional cleanup, use an `if`:

```lang
f := open(path) or return err
defer f.close()

buf := alloc(4096) or return err
success := false
defer if !success { free(buf) }

// ... do work ...
success = true
return nil
```

This is uglier than `errdefer`, but it's *explicit*. And you only pay the ugliness in the rare case where you need it.

Actually — reconsider. That `defer if !success` is a bit gross. Two cleaner options:

**Option A:** No conditional cleanup. Just `defer free(buf)` unconditionally. If `buf` escapes, don't `defer` it. This is the common case and it's fine.

**Option B:** Add `errdefer` as syntax sugar for exactly the pattern above. It's a small compiler feature (a flag on the defer entry that's checked on function return). Difficulty: **2/10**.

I'd lean toward **Option A for v1**, add `errdefer` only if real code demands it.

## Final Revised Language

You now have a language that is:

- **Manually managed** — no GC, no RC, no borrow checker.
- **Explicit** — every cleanup is a visible `defer`; every error is a visible return value.
- **Predictable** — no hidden calls, no hidden allocations, no hidden control flow.
- **Simple to compile** — AST → typed IR → codegen, with no borrow analysis, no lifetime inference, no drop glue, no monomorphization complexity (if you defer templates).
- **Familiar** — it looks like Go with sum types and manual memory. Any Go or C programmer can read it in an afternoon.

The implementation difficulty of your *entire language* is now roughly:

| Component | Difficulty |
|---|---|
| Lexer + parser | 2 |
| Type checker (no templates, no lifetimes) | 3 |
| Sum types + exhaustiveness | 3 |
| Interfaces + vtables | 3 |
| `defer` with block semantics | 3 |
| C interop | 3 |
| Codegen (LLVM or custom) | 4 |
| Build system + package resolution | 3 |
| Standard library | ∞ (but that's true of every language) |

**Total language complexity: ~3/10.** That's a language one competent engineer can implement in a year and maintain for a decade. That's a workhorse.

## A Closing Thought

You've now rejected: GC, RC, borrow checking, move semantics, destructors, exceptions, inheritance, operator overloading, implicit conversions, macros, reflection, async, and goroutines. What's left is a language where **everything the program does is visible in the source**. That's not a weakness — that's the entire point. C, Go, and Zig all succeed *because* of what they refuse to add, not despite it.

Ship this. Then add templates in v2 when real users complain. Not before.

===============================

I am already good at implementing parsing. I want a detailed plan on how to implement typechecking and conversion from AST to IR in my language, that is where I have no big experience and majority of learning materials skip or oversimplify it. It is better to restrict used techniques and algorithms to only those needed for language features we picked or plan to introduce (like templates), for defered features we can skip implementation for now, just plan and point to the place where to add them.

================================

# Typechecking and AST → IR: A Concrete Implementation Plan

This is the part every tutorial skips. Let me give you the actual algorithm, data structures, and order of operations for *your* language — no hand-waving.

I'll structure this as:

1. **The overall pipeline** (what phases exist, what each produces)
2. **Data structures** (the IR, the type representation, the symbol tables)
3. **Type checking algorithm** (walk order, how each construct is checked)
4. **AST → IR lowering** (how each construct becomes IR)
5. **Extension points** for templates, `defer` in loops, `errdefer`, etc.

I'll use a pseudo-Rust/pseudo-Go notation for the compiler's own code, since you're implementing it in *some* language and the ideas transfer.

---

## 1. The Pipeline

```
Source
  │
  ▼
[Lexer] ──► Tokens
  │
  ▼
[Parser] ──► AST (untyped, syntactic)
  │
  ▼
[Resolver] ──► AST + symbol tables (names bound to declarations)
  │
  ▼
[TypeChecker] ──► Typed AST (every expression has a type, every name resolved)
  │
  ▼
[Lowering] ──► IR (flat, typed, SSA-ish or three-address)
  │
  ▼
[Optimizer] ──► IR' (optional, can be a no-op initially)
  │
  ▼
[Codegen] ──► Machine code / LLVM IR / C
```

**Key decision:** Do you typecheck the AST directly, or lower to an untyped IR first and typecheck *that*?

- **Typecheck the AST** (recommended for you): Simpler, fewer data structures, error messages point at source. Downside: typechecking logic is entangled with AST shape.
- **Lower to untyped IR, then typecheck** (Rust, Swift): Cleaner separation, but requires a whole extra IR and a way to map errors back to source. Overkill for your feature set.

**Do this:** Typecheck the AST, produce a *typed AST*, then lower the typed AST to IR. The typed AST is a distinct set of types from the untyped AST — this is important, because it lets you use the type system of your compiler's implementation language to guarantee "this IR node is only built from typed AST nodes."

```
Untyped AST  ──[Resolver]──►  Resolved AST
Resolved AST ──[TypeChecker]──► Typed AST
Typed AST    ──[Lowering]──►  IR
```

Three AST types is a lot of boilerplate. Two is enough:

```
Untyped AST  ──[Resolver + TypeChecker]──►  Typed AST
Typed AST    ──[Lowering]──►  IR
```

The Resolver and TypeChecker can be the same pass. Actually, they *should* be separate passes because name resolution is a prerequisite for type checking (you can't type `foo(x)` until you know what `foo` refers to). But they operate on the same AST data structure, mutating it in place (filling in `resolved_symbol` fields and `type` fields).

So:

```
Untyped AST
  │
  ▼ [Pass 1: Resolver] — fills in symbol references, builds symbol tables
Resolved AST (types still unknown)
  │
  ▼ [Pass 2: TypeChecker] — fills in types, checks rules
Typed AST
  │
  ▼ [Pass 3: Lowering] — produces IR
IR
```

This is the classic three-pass frontend. It's what Go, Swift, and most production compilers do. Don't try to fuse them; you'll regret it.

---

## 2. Data Structures

### 2.1 Symbol Tables

You need two kinds:

- **Package scope**: all top-level declarations in a package (functions, types, constants, globals).
- **Lexical scopes**: locals inside a function, nested blocks.

```rust
struct Symbol {
    name: String,
    kind: SymbolKind,       // Func, Type, Const, Global, Local, Param
    ty: TypeId,             // filled in by typechecker
    decl_span: Span,        // for error messages
    visibility: Visibility, // Pub or Private
}

enum SymbolKind {
    Func { params: Vec<TypeId>, results: Vec<TypeId> },
    Type { def: TypeDefId },
    Const { value: ConstValue },
    Global { mutable: bool },
    Local { mutable: bool, index: u32 },  // index into stack frame
    Param { index: u32 },
}

struct Scope {
    symbols: HashMap<String, SymbolId>,
    parent: Option<ScopeId>,
    // for defer: list of deferred expressions registered in this scope
    defers: Vec<ExprId>,
}

struct SymbolTable {
    symbols: Arena<Symbol>,   // SymbolId -> Symbol
    scopes: Arena<Scope>,     // ScopeId -> Scope
    package_scope: ScopeId,
}
```

**Arena-based IDs, not `Box`/`Rc`.** This is important. Use integer IDs everywhere. It makes the compiler fast, avoids borrow-checker fights in your implementation language, and makes serialization/incremental compilation possible later.

```rust
struct Arena<T> { items: Vec<T> }
type SymbolId = u32;
type ScopeId = u32;
type TypeId = u32;
type ExprId = u32;
type StmtId = u32;
type FuncId = u32;
```

### 2.2 Type Representation

You have these type constructors:

- Primitives: `i8..i64`, `u8..u64`, `f32`, `f64`, `bool`, `byte`
- `string` (a built-in: pointer + length)
- Pointers: `*T`, `?*T`
- Slices: `[]T` (pointer + length)
- Fixed arrays: `[N]T`
- Structs: named, with fields
- Enums: named, with variants, each variant has a payload
- Functions: `fn(A, B) -> (C, D)`
- Interfaces: named, with a method set
- Templates (v2): `T` type parameters — but only inside a template definition

```rust
enum TypeKind {
    // Primitives
    I8, I16, I32, I64,
    U8, U16, U32, U64,
    F32, F64,
    Bool,
    String,

    // Compound
    Pointer { to: TypeId, nullable: bool },
    Slice { elem: TypeId },
    Array { elem: TypeId, len: u64 },

    // Named
    Struct { def: StructId },
    Enum { def: EnumId },
    Interface { def: InterfaceId },

    // Functions
    Func { params: Vec<TypeId>, results: Vec<TypeId> },

    // Template parameter (only valid inside a template body)
    Param { index: u32, name: String },

    // Errors: poison type, prevents cascading errors
    Error,
}

struct TypeTable {
    types: Arena<TypeKind>,
    // interning: (TypeKind) -> TypeId, so equal types have the same ID
    intern: HashMap<TypeKindKey, TypeId>,
}
```

**Intern your types.** Two `*i32` types should have the same `TypeId`. This makes type equality a single integer comparison, makes hashing trivial, and makes template instantiation cacheable.

### 2.3 The Typed AST

After typechecking, every expression node has a `ty: TypeId`. After resolution, every identifier has a `symbol: SymbolId`. The cleanest way is to make the AST generic over a "phase marker":

```rust
// Untyped
struct Expr { kind: ExprKind, span: Span }
enum ExprKind {
    IntLit(i64),
    Ident(String),                    // unresolved name
    Call(Box<Expr>, Vec<Expr>),
    Binary(BinOp, Box<Expr>, Box<Expr>),
    // ...
}

// Typed
struct TypedExpr { kind: TypedExprKind, ty: TypeId, span: Span }
enum TypedExprKind {
    IntLit(i64),
    LocalRef { index: u32 },           // resolved to a local slot
    GlobalRef { symbol: SymbolId },
    Call { callee: Callee, args: Vec<TypedExpr> },
    Binary(BinOp, Box<TypedExpr>, Box<TypedExpr>),
    // ...
}
```

You *can* use a single AST with `Option<TypeId>` fields, but I recommend two distinct types. The compiler's type system will then catch "I forgot to typecheck this node before lowering it."

### 2.4 The IR

For your language, a **three-address code (TAC)** IR is the right choice. Not SSA (too much machinery for v1), not a tree (too close to AST). TAC is:

- Flat list of instructions.
- Each instruction has at most one operator and up to 3 operands.
- Values are named (virtual registers).
- Control flow is explicit labels and jumps.

```rust
struct IrFunction {
    name: String,
    params: Vec<IrValue>,
    results: Vec<IrValue>,
    locals: Vec<IrLocal>,       // stack slots for addressable values
    blocks: Vec<IrBlock>,       // basic blocks
    entry: BlockId,
}

struct IrBlock {
    id: BlockId,
    phis: Vec<IrPhi>,           // only if you do SSA; skip for v1
    instrs: Vec<IrInstr>,
    terminator: IrTerminator,
}

enum IrInstr {
    // dst = op a, b
    BinOp { dst: IrValue, op: BinOp, a: IrValue, b: IrValue, ty: TypeId },
    UnOp  { dst: IrValue, op: UnOp, a: IrValue, ty: TypeId },

    // dst = a
    Copy { dst: IrValue, src: IrValue, ty: TypeId },

    // dst = *ptr  /  *ptr = src
    Load  { dst: IrValue, ptr: IrValue, ty: TypeId },
    Store { ptr: IrValue, src: IrValue, ty: TypeId },

    // dst = &local
    AddrOf { dst: IrValue, local: LocalId },

    // dst = call f(args)  — multiple results
    Call { dsts: Vec<IrValue>, callee: IrCallee, args: Vec<IrValue> },

    // Interface dispatch
    CallVirtual { dsts: Vec<IrValue>, iface: IrValue, method: MethodId, args: Vec<IrValue> },
}

enum IrTerminator {
    Return(Vec<IrValue>),
    Jump(BlockId),
    Branch { cond: IrValue, then_bb: BlockId, else_bb: BlockId },
    Unreachable,
}

enum IrValue {
    Reg(u32),        // virtual register
    ConstInt(i64, TypeId),
    ConstFloat(f64, TypeId),
    ConstBool(bool),
    ConstString(StringId),
    LocalAddr(LocalId),  // address of a stack local
}

enum IrCallee {
    Direct(FuncId),
    Indirect(IrValue),  // function pointer
    Extern(String),     // C function
}
```

**Two important design choices:**

1. **Values vs. locals.** A "value" is an SSA-like virtual register (single assignment, can be optimized). A "local" is a stack slot (addressable, mutable). In TAC without SSA, you can use one namespace and allow multiple assignments. I recommend: **locals for anything addressable, virtual registers for intermediate results.** This avoids a lot of phi-node pain.

2. **Basic blocks with explicit terminators.** Every block ends with exactly one terminator. This makes control flow explicit and CFG construction trivial.

---

## 3. The Resolver (Pass 1)

**Input:** Untyped AST + package graph.
**Output:** Resolved AST + symbol tables.
**Errors:** Undefined names, duplicate declarations, visibility violations.

### Algorithm

1. **Collect top-level declarations per package.** For each package, walk its top-level items and insert them into the package scope. This is a single pass, no ordering required — top-level declarations are order-independent in your language (unlike C).

2. **For each function body, walk the AST in lexical order**, maintaining a stack of scopes. On entering a block, push a scope. On leaving, pop it.

3. **On `let x = ...`:** insert `x` into the current scope. If `x` already exists in the current scope, error. (Shadowing in an outer scope is allowed; shadowing in the same scope is not.)

4. **On identifier use:** walk up the scope stack, find the nearest declaration. If none, check package scope. If none, check imported packages (respecting `pub`). If still none, error.

5. **On function call `f(args)`:** resolve `f` the same way. Record it in the AST node so the typechecker doesn't have to re-resolve.

6. **On type annotation `x: T`:** resolve `T` to a `TypeDefId`. If it's a primitive, that's fine. If it's a named type, look it up. If it's a template parameter (inside a template), record it as such.

**Key data structure:** a stack of `Scope` objects, each a `HashMap<String, SymbolId>`. Pushing/popping is O(1). Lookup is O(depth), but depth is small (typically < 10). If you want O(1), use a persistent map or a scope-indexed vector — but don't bother; O(depth) is fine.

**Interaction with `defer`:** record deferred expressions in the scope where they appear. The lowering pass will use this. The resolver doesn't need to do anything special with `defer` except resolve the expression inside it.

**Where templates go later:** template parameters live in a scope that wraps the function's body scope. The resolver inserts them as `SymbolKind::TypeParam`. Resolution of `T` finds them the same way it finds locals. This is a clean extension point.

---

## 4. The Typechecker (Pass 2)

**Input:** Resolved AST.
**Output:** Typed AST (every expression has a `TypeId`).
**Errors:** Type mismatches, arity mismatches, missing fields, non-exhaustive match, etc.

### 4.1 The Core Algorithm: Bidirectional Checking

You have two modes:

- **Synthesize (infer):** given an expression, produce its type. `synth(e) -> TypeId`
- **Check:** given an expression and an expected type, verify they match. `check(e, expected) -> ()`

The rules:

```
synth(literal)         = its natural type
synth(ident)           = symbol's type
synth(call f(args))    = results of f's type, after checking args against params
synth(binary op a b)   = op's result type, after checking a and b

check(e, expected) =
    if e is a literal or a construct with a "default" type:
        synth(e) and verify it equals expected (or is coercible)
    else:
        synth(e) and verify it equals expected
```

For your language, **coercion is almost nonexistent** (you banned implicit conversions). The only coercions you might allow:

- `T` → `?*T` (non-null pointer to nullable pointer) — implicit, safe.
- `[N]T` → `[]T` (array to slice) — implicit, safe.
- Integer literal → any integer type, if the value fits. This is important; otherwise `let x: u32 = 5` requires `5u32`.

Everything else is an error. This makes the typechecker *dramatically* simpler than languages with rich coercion rules.

### 4.2 Walk Order

Walk the AST top-down, depth-first. For each node:

1. Recursively typecheck children.
2. Apply the node's own rule.
3. Store the resulting type in the node.

For statements, typechecking is mostly about scoping and control flow:

```rust
fn check_stmt(&mut self, stmt: &Stmt) {
    match stmt {
        Stmt::Let { name, ty_annot, init } => {
            let init_ty = match ty_annot {
                Some(t) => { self.check_expr(init, t); t }
                None    => self.synth_expr(init),  // type inference for locals
            };
            self.declare_local(name, init_ty);
        }
        Stmt::Assign { lhs, rhs } => {
            let lhs_ty = self.synth_expr(lhs);
            self.check_expr(rhs, lhs_ty);
            if !self.is_assignable(lhs) { error("cannot assign to this expression") }
        }
        Stmt::Expr(e) => { self.synth_expr(e); }
        Stmt::If { cond, then, else_ } => {
            self.check_expr(cond, BOOL);
            self.check_block(then);
            if let Some(e) = else_ { self.check_block(e); }
        }
        Stmt::For { init, cond, step, body } => {
            self.push_scope();
            if let Some(i) = init { self.check_stmt(i); }
            if let Some(c) = cond { self.check_expr(c, BOOL); }
            if let Some(s) = step { self.check_stmt(s); }
            self.check_block(body);
            self.pop_scope();
        }
        Stmt::Return(exprs) => {
            let expected = self.current_function_results();
            self.check_return(exprs, expected);
        }
        Stmt::Defer(e) => {
            self.synth_expr(e);  // just check it's valid
            self.current_scope_defers.push(e.id);
        }
        Stmt::Match { scrutinee, arms } => { self.check_match(scrutinee, arms); }
        Stmt::Block(b) => { self.push_scope(); self.check_block(b); self.pop_scope(); }
    }
}
```

### 4.3 Type Inference for `let` Without Annotation

You said no complex type systems. Good. So inference is **local only**:

- `let x = expr` — infer `x`'s type from `expr`'s synthesized type.
- No Hindley-Milner.
- No inference across function boundaries.
- No inference for function parameters or return types.

This means `synth_expr` never needs to "guess" — it always has enough information. The only place inference happens is `let` without annotation, and it's a one-liner.

If `expr` is a composite literal like `{1, 2, 3}`, you can't infer it without a target type. So:

- `let x = [1, 2, 3]` — infer `[3]i32` (array of 3 i32, using default int type).
- `let x: [3]u8 = [1, 2, 3]` — check each element against `u8`.
- `let x = {a: 1, b: 2}` — error: struct literals need a type annotation.

This is a deliberate restriction. It avoids the "inference propagates backward" problem.

### 4.4 Function Calls

```rust
fn synth_call(&mut self, callee: &Expr, args: &[Expr]) -> TypeId {
    let callee_ty = self.synth_expr(callee);
    match self.types.get(callee_ty) {
        TypeKind::Func { params, results } => {
            if args.len() != params.len() {
                error(format!("expected {} args, got {}", params.len(), args.len()));
                return ERROR_TYPE;
            }
            for (arg, param_ty) in args.iter().zip(params) {
                self.check_expr(arg, *param_ty);
            }
            // Multiple return values: return a "tuple" type
            self.tuple_type(results.clone())
        }
        TypeKind::Error => ERROR_TYPE,
        _ => { error("called expression is not a function"); ERROR_TYPE }
    }
}
```

**Multiple return values:** model them as a tuple type. When you have `fn f() -> (i32, Error?)`, its `results` field is `vec![I32, OPTIONAL_ERROR]`. A call to `f` synthesizes a tuple type. The caller can:

- `let (a, b) = f()` — destructure.
- `let x = f()` — if you allow tuple-typed variables.
- Use `f().0` — if you have tuple indexing. (Maybe skip this.)

For simplicity, I recommend **no first-class tuples**. Multiple returns are only usable via destructuring `let (a, b) = f()` and via `return (x, y)`. Internally you can represent the tuple type, but don't let users store it in a variable or pass it around. This keeps the type system and IR simple.

### 4.5 Structs

```rust
fn check_struct_literal(&mut self, name: TypeId, fields: &[(String, Expr)]) -> TypeId {
    let def = self.struct_def(name);
    let mut seen = HashSet::new();
    for (fname, fexpr) in fields {
        let field_ty = match def.fields.get(fname) {
            Some(t) => *t,
            None => { error(format!("no field {}", fname)); continue; }
        };
        if !seen.insert(fname) { error(format!("duplicate field {}", fname)); }
        self.check_expr(fexpr, field_ty);
    }
    for f in &def.fields {
        if !seen.contains(&f.name) {
            error(format!("missing field {}", f.name));
        }
    }
    name
}
```

Field access `x.f`:

```rust
fn synth_field_access(&mut self, obj: &Expr, field: &str) -> TypeId {
    let obj_ty = self.synth_expr(obj);
    match self.types.get(obj_ty) {
        TypeKind::Struct { def } => {
            let s = self.struct_defs.get(def);
            match s.fields.get(field) {
                Some(t) => *t,
                None => { error(format!("no field {}", field)); ERROR_TYPE }
            }
        }
        TypeKind::Pointer { to, .. } => {
            // auto-deref: (*p).f
            self.synth_field_access_deref(*to, field)
        }
        _ => { error("field access on non-struct"); ERROR_TYPE }
    }
}
```

**Auto-deref for pointers** is a convenience. Decide whether you want it. Go does it; C doesn't. I'd allow it for `*T` (you write `p.field` instead of `(*p).field`), but require explicit `*` for `?*T` after a null check.

### 4.6 Enums and `match`

An enum has variants. Each variant has a payload (possibly empty, possibly multiple fields).

```rust
struct EnumDef {
    name: String,
    variants: Vec<Variant>,
}
struct Variant {
    name: String,
    fields: Vec<TypeId>,   // empty for unit variants
}
```

The type of an enum is the enum itself. Each variant is *not* a separate type (unlike Rust, where `Option::Some` is a function returning `Option<T>`). This is simpler: `Token::Ident` is just a tag, and `Token::Ident("foo")` is a construction expression that produces a `Token`.

```rust
fn check_enum_construct(&mut self, enum_ty: TypeId, variant: &str, args: &[Expr]) -> TypeId {
    let def = self.enum_def(enum_ty);
    let v = match def.variants.iter().find(|v| v.name == variant) {
        Some(v) => v,
        None => { error("no such variant"); return ERROR_TYPE; }
    };
    if args.len() != v.fields.len() {
        error("wrong number of fields in variant");
    }
    for (arg, field_ty) in args.iter().zip(&v.fields) {
        self.check_expr(arg, *field_ty);
    }
    enum_ty
}
```

`match`:

```rust
fn check_match(&mut self, scrutinee: &Expr, arms: &[MatchArm]) -> TypeId {
    let scrut_ty = self.synth_expr(scrutinee);
    let enum_def = match self.types.get(scrut_ty) {
        TypeKind::Enum { def } => self.enum_defs.get(def),
        _ => { error("match on non-enum"); return ERROR_TYPE; }
    };

    let mut covered: HashSet<VariantId> = HashSet::new();
    let mut result_ty: Option<TypeId> = None;

    for arm in arms {
        let (variant, bindings) = parse_pattern(&arm.pattern);
        if !covered.insert(variant) {
            error("duplicate arm");
        }
        let v = enum_def.variants.get(variant);
        if bindings.len() != v.fields.len() {
            error("pattern arity mismatch");
        }
        self.push_scope();
        for (name, field_ty) in bindings.iter().zip(&v.fields) {
            self.declare_local(name, *field_ty);
        }
        let arm_ty = self.check_expr(&arm.body, result_ty.unwrap_or(UNKNOWN));
        match result_ty {
            None => result_ty = Some(arm_ty),
            Some(t) if t == arm_ty => {}
            Some(t) => error("match arms have different types"),
        }
        self.pop_scope();
    }

    if covered.len() != enum_def.variants.len() {
        error("non-exhaustive match");
    }

    result_ty.unwrap_or(UNIT)
}
```

**Exhaustiveness** is a set comparison: every variant must be covered. This is trivial for flat enums. For nested patterns (`Some(Ok(x))`), you need a recursive algorithm (Maranget's), but you banned nested patterns, so set comparison is enough.

### 4.7 Interfaces

An interface declares a method set. A type satisfies an interface if its method set is a superset of the interface's.

```rust
fn check_interface_assignment(&mut self, value_ty: TypeId, iface_ty: TypeId) -> bool {
    let iface = self.interface_def(iface_ty);
    let methods = self.method_set_of(value_ty);
    for required in &iface.methods {
        match methods.get(&required.name) {
            Some(m) if m.signature == required.signature => {}
            _ => return false,
        }
    }
    true
}

fn method_set_of(&self, ty: TypeId) -> HashMap<String, MethodSig> {
    match self.types.get(ty) {
        TypeKind::Struct { def } => self.struct_defs.get(def).methods.clone(),
        TypeKind::Pointer { to, .. } => self.method_set_of(*to),  // maybe
        TypeKind::Interface { def } => self.interface_defs.get(def).methods.clone(),
        _ => HashMap::new(),
    }
}
```

**Method set computation is a fixpoint** if you allow interface embedding. If you ban embedding (recommended), it's a direct lookup.

**Where does the vtable come from?** The typechecker records, for each `(concrete_ty, iface_ty)` pair that's assigned, that a vtable is needed. The lowering pass generates it. More on this in §5.7.

### 4.8 Error Recovery

**Do not stop at the first error.** Return a special `ERROR_TYPE` from any failed synthesis. `ERROR_TYPE` is compatible with *everything* (it unifies with any type). This prevents cascading errors: one mistake produces one diagnostic, not fifty.

```rust
fn types_compatible(&self, a: TypeId, b: TypeId) -> bool {
    if a == ERROR_TYPE || b == ERROR_TYPE { return true; }
    a == b  // after interning, equality is pointer equality
}
```

Report errors via a `Diagnostics` sink that collects them. Don't panic. Don't `return Err` from every function — that's what `ERROR_TYPE` is for.

---

## 5. AST → IR Lowering (Pass 3)

**Input:** Typed AST.
**Output:** IR.
**Errors:** None (all errors caught in typechecking). Lowering should be infallible.

### 5.1 The Core Idea

Lowering is a recursive tree walk that **emits IR instructions into the current basic block** and **returns an `IrValue`** representing the result of the expression.

```rust
struct Lowering<'a> {
    func: &'a mut IrFunction,
    current_block: BlockId,
    // map from local symbol -> LocalId
    locals: HashMap<SymbolId, LocalId>,
    // map from scope to list of defers
    defer_stack: Vec<Vec<DeferredExpr>>,
}

fn lower_expr(&mut self, expr: &TypedExpr) -> IrValue {
    match &expr.kind {
        TypedExprKind::IntLit(n) => self.const_int(*n, expr.ty),
        TypedExprKind::LocalRef { index } => {
            // load from local slot
            let dst = self.new_reg(expr.ty);
            self.emit(IrInstr::Load { dst, ptr: self.local_addr(*index), ty: expr.ty });
            dst
        }
        TypedExprKind::Binary(op, a, b) => {
            let av = self.lower_expr(a);
            let bv = self.lower_expr(b);
            let dst = self.new_reg(expr.ty);
            self.emit(IrInstr::BinOp { dst, op: *op, a: av, b: bv, ty: expr.ty });
            dst
        }
        TypedExprKind::Call { callee, args } => {
            let arg_vals: Vec<_> = args.iter().map(|a| self.lower_expr(a)).collect();
            let result_tys = self.call_result_types(callee);
            let dsts: Vec<_> = result_tys.iter().map(|t| self.new_reg(*t)).collect();
            self.emit(IrInstr::Call { dsts: dsts.clone(), callee: self.lower_callee(callee), args: arg_vals });
            if dsts.len() == 1 { dsts[0] } else { /* tuple — see below */ }
        }
        // ...
    }
}
```

### 5.2 Multiple Return Values in IR

A function `fn f() -> (i32, Error?)` has two result registers. A call to it produces two `IrValue`s. But `lower_expr` returns a single `IrValue`.

**Solution:** Don't lower a multi-return call as an expression. Lower it as a *statement*:

```rust
Stmt::LetTuple { names, init } => {
    let (dsts, call_instr) = self.lower_call_multi(init);
    // dsts is Vec<IrValue>
    // bind names to dsts
}
```

So `let (a, b) = f()` is a distinct statement form. `f()` alone as an expression is only allowed if `f` returns exactly one value. This is a restriction, but a clean one.

**Alternative:** Represent a tuple as a pointer to a stack-allocated struct. `f` writes its results to a caller-provided pointer. This is what C ABI does for large structs. It's more uniform but less efficient. I'd go with multi-register results and a special statement form.

### 5.3 Control Flow

```rust
Stmt::If { cond, then, else_ } => {
    let cond_v = self.lower_expr(cond);
    let then_bb = self.new_block();
    let else_bb = self.new_block();
    let merge_bb = self.new_block();

    self.terminate(IrTerminator::Branch { cond: cond_v, then_bb, else_bb });

    self.current_block = then_bb;
    self.lower_block(then);
    self.terminate(IrTerminator::Jump(merge_bb));

    self.current_block = else_bb;
    if let Some(e) = else_ { self.lower_block(e); }
    self.terminate(IrTerminator::Jump(merge_bb));

    self.current_block = merge_bb;
}
```

Every construct that can branch creates new basic blocks and terminates the current one. This is mechanical and hard to get wrong if you follow the pattern:

1. Evaluate operands (may emit instructions into the current block).
2. Terminate the current block with a branch.
3. Create new blocks for each arm.
4. Lower each arm in its own block, terminating with a jump to the merge block.
5. Continue in the merge block.

**Loops** follow the same pattern with `for`:

```rust
Stmt::For { init, cond, step, body } => {
    self.lower_opt_stmt(init);
    let header_bb = self.new_block();
    let body_bb = self.new_block();
    let exit_bb = self.new_block();

    self.terminate(IrTerminator::Jump(header_bb));

    self.current_block = header_bb;
    let cond_v = match cond {
        Some(c) => self.lower_expr(c),
        None => self.const_bool(true),
    };
    self.terminate(IrTerminator::Branch { cond: cond_v, then_bb: body_bb, else_bb: exit_bb });

    self.current_block = body_bb;
    // push a defer scope for the loop body
    self.push_defer_scope();
    self.lower_block(body);
    // pop defer scope, emitting defers at end of body
    self.pop_defer_scope();
    self.lower_opt_stmt(step);
    self.terminate(IrTerminator::Jump(header_bb));

    self.current_block = exit_bb;
}
```

### 5.4 `defer` with Block Semantics

This is the trickiest part of lowering. The rule:

- `defer e` registers `e` in the current scope's defer list.
- When the scope exits (normally, via `break`, `continue`, or `return`), all registered defers run in reverse order.

**Implementation strategy:** Don't try to inline defers at every exit point. Instead:

1. **Lower each `defer` body into a separate "cleanup block"** during the block's lowering. The cleanup block contains the lowered defer expression followed by a jump to the scope's exit block.
2. **Record the cleanup block ID** in the scope's defer list.
3. **At every exit from the scope** (end of block, `break`, `continue`, `return`), emit a jump to the *chain* of cleanup blocks.

Actually, this is over-engineered for v1. Simpler approach:

**Strategy A: Duplicate defer bodies at each exit.** When lowering a block that contains defers, at every exit point (normal end, `break`, `continue`, `return`), lower the deferred expressions in reverse order into the current block before emitting the exit terminator.

This causes code duplication if there are many exits, but it's simple and correct. For v1, do this.

**Strategy B: Cleanup chain.** More efficient, more complex. Skip for now.

Let me sketch Strategy A:

```rust
struct DeferScope {
    defers: Vec<TypedExpr>,   // in registration order
}

fn lower_block(&mut self, block: &TypedBlock) {
    self.push_defer_scope();
    for stmt in &block.stmts {
        match stmt {
            Stmt::Defer(e) => self.current_defer_scope().defers.push(e.clone()),
            Stmt::Return(exprs) => {
                // lower return values first
                let vals = exprs.iter().map(|e| self.lower_expr(e)).collect();
                // then run defers in reverse order
                self.emit_defers_for_current_scope_and_parents();
                self.terminate(IrTerminator::Return(vals));
            }
            Stmt::Break => {
                self.emit_defers_up_to_loop_scope();
                self.terminate(IrTerminator::Jump(self.loop_exit_bb));
            }
            _ => self.lower_stmt(stmt),
        }
    }
    // normal fall-through: emit defers, then continue
    self.emit_defers_for_current_scope();
    self.pop_defer_scope();
}

fn emit_defers_for_current_scope(&mut self) {
    let defers = self.current_defer_scope().defers.clone();
    for d in defers.iter().rev() {
        self.lower_expr(d);   // result discarded
    }
}
```

**Important detail:** `defer` bodies are lowered *at the exit point*, not at the `defer` statement. This means the defer expression sees the state of variables *at exit*, not at registration. This is Go's semantics and Zig's, and it's what users expect. But it also means you can't `defer` a call whose arguments depend on values that are moved/changed — well, actually you can, because you don't have move semantics. Good.

**`return` runs all defers from inner scopes up to the function scope.** This is the "run defers in all enclosing scopes" rule. Implement by walking the defer scope stack from top to bottom (excluding the function scope itself, which is popped at function end).

**`break`/`continue` run defers from inner scopes up to (but not including) the loop's scope.** The loop body's defers run, but the loop header's don't.

### 5.5 Match Lowering

`match` lowers to a chain of branches on the tag:

```rust
Stmt::Match { scrutinee, arms } => {
    let scrut_v = self.lower_expr(scrutinee);
    let scrut_local = self.new_local(scrutinee.ty);
    self.emit(IrInstr::Store { ptr: self.local_addr(scrut_local), src: scrut_v, ty: scrutinee.ty });

    let merge_bb = self.new_block();
    let mut arm_blocks = Vec::new();
    for arm in arms {
        let bb = self.new_block();
        arm_blocks.push((bb, arm));
    }

    // dispatch: load tag, compare, jump
    for (i, (bb, arm)) in arm_blocks.iter().enumerate() {
        let tag = self.emit_load_tag(scrut_local);
        let expected_tag = self.const_int(arm.variant_index as i64, U32);
        let cmp = self.emit_binop(Eq, tag, expected_tag, BOOL);
        let next_bb = if i + 1 < arm_blocks.len() { arm_blocks[i+1].0 } else { self.unreachable_bb() };
        let this_bb = self.new_block();
        self.terminate(IrTerminator::Branch { cond: cmp, then_bb: *bb, else_bb: next_bb });
        self.current_block = this_bb;
        // ... this is getting verbose; in practice you'd generate a switch
    }

    // lower each arm
    for (bb, arm) in arm_blocks {
        self.current_block = bb;
        // bind payload fields to locals by loading from scrut_local
        let payload_base = self.emit_payload_ptr(scrut_local);
        for (i, binding) in arm.bindings.iter().enumerate() {
            let field_ptr = self.emit_ptr_offset(payload_base, i * field_size);
            let val = self.emit_load(field_ptr, binding.ty);
            let local = self.new_local(binding.ty);
            self.emit(Store { ptr: self.local_addr(local), src: val, ty: binding.ty });
            self.bind_local(binding.symbol, local);
        }
        let result_v = self.lower_expr(&arm.body);
        self.emit(Store { ptr: self.result_local, src: result_v, ty: arm.body.ty });
        self.terminate(IrTerminator::Jump(merge_bb));
    }

    self.current_block = merge_bb;
    // load result from result_local
}
```

**Enum representation:** tag (u32) + payload (union of variant field types, sized to the largest). The tag is stored first, then the payload, aligned appropriately.

```
Enum Token:
  tag: u32
  payload: union {
      Ident: { ptr: *u8, len: usize },
      Int: { value: i64 },
      Plus: {},
      Minus: {},
  }
```

Total size = `align_to(max(sizeof(payload), align_of(payload)) + sizeof(tag), max_align)`.

### 5.6 Interface Lowering

An interface value is a **fat pointer**: `{ data: *void, vtable: *VTable }`.

- `data` points to the concrete value (heap-allocated if the value doesn't fit in a pointer, or if you require it).
- `vtable` points to a static array of function pointers, one per method in the interface.

Assigning a concrete value to an interface variable:

1. **Box the value** if it's not already a pointer. (This is an allocation — you should decide whether it's implicit. For a manual memory language, it *must* be explicit. So: `let r: Reader = Reader(&my_file)` where `&my_file` is explicit.)
2. **Look up the vtable** for `(concrete_ty, iface_ty)`. Generate it if not already generated.
3. **Build the fat pointer** `{ data, vtable }`.

Calling a method on an interface:

1. Load `vtable` from the fat pointer.
2. Load the function pointer at the method's index.
3. Call it with `data` as the first argument.

**Lowering:**

```rust
TypedExprKind::InterfaceCall { iface, method, args } => {
    let iface_v = self.lower_expr(iface);   // fat pointer value
    let data = self.emit_extract_field(iface_v, 0, PTR_VOID);
    let vtable = self.emit_extract_field(iface_v, 1, PTR_VTABLE);
    let method_slot = self.emit_vtable_load(vtable, method.index);
    let arg_vals: Vec<_> = std::iter::once(data).chain(args.iter().map(|a| self.lower_expr(a))).collect();
    let dst = self.new_reg(method.result_ty);
    self.emit(IrInstr::Call { dsts: vec![dst], callee: IrCallee::Indirect(method_slot), args: arg_vals });
    dst
}
```

**Vtable generation:** For each `(concrete_ty, iface_ty)` pair used in the program, emit a static array of function pointers at the top level. The typechecker collects these pairs; the lowering pass emits the tables.

### 5.7 C Interop Lowering

`extern fn` declarations lower to `IrCallee::Extern("malloc")`. The call is emitted as `Call` with the extern callee. The codegen stage resolves the symbol at link time.

`export fn` definitions lower normally, but the codegen stage emits them with C ABI (name mangling disabled, calling convention set).

**ABI differences:** Your language's calling convention vs. C's. For simple cases (scalars, pointers), they're the same. For structs by value, they differ. For v1, restrict `extern`/`export` to functions whose parameters and returns are scalars, pointers, or slices (which are just `{ptr, len}` structs). Document this restriction.

---

## 6. Extension Points for Deferred Features

### 6.1 Templates

**Where they plug in:**

- **Resolver:** when entering a template definition, push a scope containing the template parameters as `SymbolKind::TypeParam`. When resolving a type name, if it matches a type parameter, record it as a `TypeParamRef` rather than a `TypeDef`.
- **Typechecker:** a template body is typechecked *once*, in a special mode where type parameters are opaque types. Each use of `T` synthesizes the type `Param { index, name }`. Operations on `T` are checked against the constraint (or, if no constraints, deferred to instantiation time).
- **Instantiation:** when you see `map[i32, i32](...)`, create a substitution `{T -> i32, U -> i32}`, clone the template's typed AST with the substitution applied, and typecheck the clone *with concrete types*. Cache instantiations by `(template_id, substitution)`.
- **Lowering:** each instantiation is lowered as a separate function. Mangle the name: `map$i32$i32`. No special handling needed.

**The hard part:** deciding when to check the template body. If you check at definition time, you need constraints. If you check at instantiation time, you don't need constraints but errors are bad. I recommend: **check at definition time with no constraints**, treating `T` as an opaque type that supports only "copy", "compare for equality", and "pass to a function". This catches most errors early and is simple. It's a restriction (no arithmetic on `T`), but a reasonable one.

### 6.2 `errdefer`

**Where it plugs in:**

- **Resolver:** same as `defer`.
- **Typechecker:** same as `defer`.
- **Lowering:** `errdefer` bodies are emitted only in the "error return" path. Track a flag `is_error_path` during block lowering. When lowering a `return` statement, if the returned error value is non-nil, emit `errdefer` bodies; otherwise, skip them.

Simpler: require that `return err` (where `err` is non-nil) and `return nil` are syntactically distinguishable. Then `errdefer` bodies are emitted before `return err` but not before `return nil`.

### 6.3 More complex pattern matching

**Where it plugs in:** the `check_match` function and the match lowering. For nested patterns, replace the set-comparison exhaustiveness check with Maranget's algorithm, and replace the flat tag switch with a decision tree. This is a localized change.

### 6.4 Method dispatch optimization

**Where it plugs in:** the lowering of `InterfaceCall`. If the concrete type is statically known (devirtualization), emit a direct call. This is an optimization pass over the IR, not a frontend change.

---

## 7. Suggested Implementation Order

1. **Resolver + symbol tables** — get name resolution working on a subset (functions, locals, calls).
2. **Typechecker for primitives and functions** — no structs, no enums, no interfaces.
3. **Lowering for primitives and functions** — get "hello world" through the pipeline.
4. **Pointers and slices** — add to typechecker and lowering.
5. **Structs** — add to typechecker (field access, literals) and lowering.
6. **Enums and match** — add to typechecker (exhaustiveness) and lowering (tag dispatch).
7. **Interfaces** — add to typechecker (method sets, satisfaction) and lowering (fat pointers, vtables).
8. **`defer`** — add to resolver, typechecker, and lowering (block semantics).
9. **C interop** — add `extern`/`export` to resolver and lowering, and to codegen.
10. **Templates** — add to all three passes. This is v2 work.

At each step, write tests. Golden tests (source + expected IR) are invaluable for lowering. Type-checking tests should be "this program typechecks" or "this program produces this exact error."

---

## 8. Common Pitfalls

1. **Don't typecheck during parsing.** The parser should produce an AST with no type information. Mixing them makes both harder to test.

2. **Don't use `Rc<RefCell<...>>` everywhere.** Arena + integer IDs. Your future self will thank you.

3. **Don't try to be clever with error recovery.** `ERROR_TYPE` and continue. One diagnostic per mistake.

4. **Don't lower before typechecking.** The typed AST is a *different* data structure. Lowering from the untyped AST means re-doing type-directed decisions.

5. **Don't use SSA in v1.** TAC with explicit locals is fine. SSA is an optimization, not a requirement. Add it later as an IR-to-IR pass if you need it.

6. **Don't inline defers.** Duplicate them at exit points. It's less efficient but correct and simple. Optimize later.

7. **Intern your types.** Equality as integer comparison is worth the small complexity.

8. **Test the resolver independently.** Name resolution bugs masquerade as type errors, which makes debugging miserable.

9. **Write the IR printer early.** You will stare at IR dumps constantly. Make them readable.

10. **Keep the typechecker and lowering in sync.** When you add a new AST node, add it to both. Use exhaustive `match` (no `_ =>` catch-all) so the compiler tells you what you forgot.

---

## 9. A Minimal End-to-End Example

Source:

```lang
fn add(a: i32, b: i32) -> i32 {
    return a + b
}

fn main() -> Error? {
    x := add(2, 3)
    if x > 4 {
        return nil
    }
    return Error.Custom(1, "too small")
}
```

After resolver: every `add`, `a`, `b`, `x`, `Error`, `Custom` is bound to a symbol.

After typechecker: every expression has a type. `add` has type `fn(i32, i32) -> i32`. `x` has type `i32`. `x > 4` has type `bool`. The function returns `Error?`.

IR for `add`:

```
func add(a: i32, b: i32) -> i32:
  entry:
    %0 = load a
    %1 = load b
    %2 = add %0, %1
    return %2
```

IR for `main`:

```
func main() -> Error?:
  entry:
    %0 = const 2 : i32
    %1 = const 3 : i32
    %2 = call add(%0, %1) -> i32
    store x, %2
    %3 = load x
    %4 = const 4 : i32
    %5 = gt %3, %4 : bool
    branch %5, then, else

  then:
    %6 = const nil : Error?
    return %6

  else:
    %7 = const 1 : i32
    %8 = const "too small" : string
    %9 = construct Error.Custom(%7, %8) : Error?
    return %9
```

This is the shape of the output you're aiming for. Simple, explicit, no magic.

---

## Summary

- **Three passes:** Resolver → Typechecker → Lowering. Don't fuse them.
- **Arena + integer IDs** everywhere. No `Rc<RefCell>`.
- **Bidirectional type checking** with `synth`/`check`. Local inference only.
- **`ERROR_TYPE`** for recovery. One error, one diagnostic.
- **TAC IR** with basic blocks and explicit terminators. No SSA in v1.
- **`defer` via duplication** at exit points. Block semantics, LIFO order.
- **Enums as tag + union payload.** Match as tag dispatch.
- **Interfaces as fat pointers.** Vtables generated per `(concrete, interface)` pair.
- **Templates** plug in as a substitution pass at the typed-AST level.
- **Test each pass independently.** Golden tests for IR, "typechecks/doesn't" for the checker.

That's the whole game. The rest is grinding through the cases.

=========================

Give me more directions on how to approach hoisting of user-defined types and constants, because constant can affect type definitions (for arrays). Give me a plan how to implement it.

=========================

# Hoisting User-Defined Types and Constants: A Concrete Plan

This is one of those problems that looks trivial until you hit the circular dependency: **arrays need constant lengths, constants can reference types, and types can reference constants.** You need a principled order of evaluation.

Let me lay out the problem precisely, then give you an algorithm.

---

## 1. The Problem

Consider these declarations (all at package scope, in any order):

```lang
const N = 16
const M = N * 2

type Buffer = struct {
    data: [M]u8
    len:  usize
}

const EMPTY = Buffer{ data: [M]u8{}, len: 0 }

type Aligned = struct {
    bytes: [align_of(Buffer)]u8
}
```

Dependencies:

- `M` depends on `N`.
- `Buffer` depends on `M` (via `[M]u8`).
- `EMPTY` depends on `Buffer` and `M`.
- `Aligned` depends on `Buffer` (via `align_of`).

If you process declarations in source order, you get it wrong. If you process in reverse, you get it wrong differently. You need **dependency-ordered evaluation**.

And there are cycles to worry about:

```lang
const A = B + 1
const B = A + 1       // cycle — error

type T = struct { x: [S]u8 }
const S = size_of(T)  // cycle — error
```

---

## 2. Key Insight: Split Types into Two Kinds

The single most important design decision here is to **split type declarations into two categories**:

1. **Structural types**: structs, enums, interfaces. Their *shape* depends only on other types, never on constants. A struct field `x: [N]u8` has a shape that depends on `N`, but a struct field `x: T` where `T` is a named type only depends on `T`'s shape.

   Actually, wait — a struct field *can* depend on a constant: `x: [N]u8`. So structural types *do* depend on constants. Hmm.

   Let me refine: **The *layout* of a struct depends on the *layout* of its fields. The layout of `[N]u8` depends on the value of `N`.** So there's no clean separation.

2. **Alias types**: `type Foo = Bar`. These are just renames.

So the "split types" idea doesn't cleanly work. Let me give you the real solution.

---

## 3. The Real Solution: A Unified Dependency Graph

Treat **all package-scope declarations** (types, constants, functions, globals) as nodes in a single directed graph. Edges represent "A depends on B". Topologically sort the graph. Evaluate in topological order. Detect cycles and report them as errors.

This is the standard approach used by every serious compiler (Go, Swift, Rust, C++ for constants). It's simple to describe and implementable in a few hundred lines.

### 3.1 Node Kinds

```rust
enum DeclId { u32 }

enum DeclKind {
    TypeAlias,       // type Foo = Bar
    Struct,          // type Foo = struct { ... }
    Enum,            // type Foo = enum { ... }
    Interface,       // type Foo = interface { ... }
    Const,           // const N = expr
    GlobalVar,       // var x: T = expr
    Function,        // fn f(...) -> ...
}
```

Note: functions are declarations too, but their *bodies* are not part of the dependency graph (only their *signatures* are). This is important — a function body can reference any declaration in the package, but that reference doesn't create a dependency edge between the function and the referenced declaration for *hoisting* purposes. The function is evaluated lazily when called.

So: **function signatures participate in hoisting; function bodies do not.**

### 3.2 Dependency Edges

For each declaration, extract the set of declarations it references. This is a **syntactic** analysis — walk the declaration's AST and collect all identifier references that resolve to package-scope symbols.

But wait — resolution requires knowing what the identifiers refer to, which requires... resolution. Circular again.

**Resolution:** Do name resolution in two phases:

**Phase A: Collect all top-level names.** Walk all package-scope declarations, collect their names into the package scope. Do NOT resolve references yet. Just collect names and their kinds (is `Foo` a type or a value?).

**Phase B: Resolve references within each declaration.** Now that all names are known, walk each declaration's body and resolve identifiers to declarations. Record the references.

This works because **the set of names is independent of the order of declarations.** You can collect them in one pass.

```rust
// Phase A: collect all top-level names
for decl in package.decls {
    match decl {
        TypeDecl { name, .. } => package_scope.insert(name, SymbolKind::Type(decl_id)),
        ConstDecl { name, .. } => package_scope.insert(name, SymbolKind::Const(decl_id)),
        FuncDecl { name, .. } => package_scope.insert(name, SymbolKind::Func(decl_id)),
        GlobalDecl { name, .. } => package_scope.insert(name, SymbolKind::Global(decl_id)),
    }
}

// Phase B: resolve references
for decl in package.decls {
    let refs = collect_refs(decl, &package_scope);
    dep_graph.add_node(decl.id);
    for r in refs {
        dep_graph.add_edge(decl.id, r.decl_id);
    }
}
```

### 3.3 What Counts as a Reference?

For each declaration kind, walk its AST and collect identifier references:

- **Type alias `type Foo = Bar`:** reference to `Bar`.
- **Struct `type Foo = struct { x: T, y: U }`:** references to `T` and `U`. If `T` is `[N]u8`, also reference to `N`.
- **Enum `type Foo = enum { A(T), B }`:** references to `T` and any other payload types.
- **Interface `type Foo = interface { fn m(x: T) -> U }`:** references to `T` and `U`.
- **Const `const N = expr`:** references to all identifiers in `expr`.
- **Global `var x: T = expr`:** reference to `T` (for the type) and all identifiers in `expr`. But globals can be lazily initialized; for hoisting purposes, only the *type* matters. Let's say globals participate in the graph only via their type.
- **Function `fn f(a: T, b: U) -> V`:** references to `T`, `U`, `V`. Body is NOT walked.

**Important:** local variables inside function bodies and constant expressions create their own scopes; don't collect them as package-level references. Use the resolver's scope stack.

### 3.4 Topological Sort with Cycle Detection

Use **Tarjan's strongly-connected components (SCC) algorithm** or **Kahn's algorithm**. Tarjan's is better because it gives you the SCCs, and each SCC of size > 1 (or a self-loop) is a cycle you can report precisely.

```rust
fn topological_sort(graph: &DepGraph) -> Result<Vec<DeclId>, Vec<Cycle>> {
    let mut order = Vec::new();
    let mut state = vec![Unvisited; graph.node_count()];
    let mut stack = Vec::new();

    fn visit(n: DeclId, graph: &DepGraph, state: &mut Vec<State>, stack: &mut Vec<DeclId>, order: &mut Vec<DeclId>) -> Result<(), Cycle> {
        match state[n] {
            Done => return Ok(()),
            InProgress => {
                // Found a cycle: extract it from the stack
                let start = stack.iter().position(|&x| x == n).unwrap();
                return Err(Cycle(stack[start..].to_vec()));
            }
            Unvisited => {}
        }
        state[n] = InProgress;
        stack.push(n);
        for &dep in graph.deps_of(n) {
            visit(dep, graph, state, stack, order)?;
        }
        stack.pop();
        state[n] = Done;
        order.push(n);
        Ok(())
    }

    for n in graph.nodes() {
        visit(n, graph, &mut state, &mut stack, &mut order)?;
    }
    Ok(order)
}
```

The resulting `order` is a topological order: **dependencies come before dependents**. This is exactly what you want.

**Cycle error message:** the cycle path tells the user exactly what's wrong. "Const `A` depends on `B`, `B` depends on `A`."

### 3.5 The Catch: Types Must Be Defined Before Use in Signatures

There's one more subtlety. Consider:

```lang
fn f(x: T) -> U { ... }
type T = struct { ... }
type U = struct { ... }
```

Function `f`'s signature references `T` and `U`. These must be *defined* (layout computed) before `f`'s signature can be typechecked. The topological sort handles this: `T` and `U` come before `f`.

But there's a bootstrapping issue: to compute the layout of `T`, you need the layout of its fields, which may include other types. The topological sort ensures fields are laid out before the struct itself.

**Rule:** A type declaration's *layout* can be computed once all types it references have their layouts computed. This is exactly what the topological order gives you.

### 3.6 What Does "Evaluating" a Declaration Mean?

In topological order, process each declaration:

- **Type alias `type Foo = Bar`:** record `Foo` as an alias for `Bar`'s type. No layout computation needed.
- **Struct `type Foo = struct { ... }`:** compute layout. Each field's type must already be resolved (it is, by topological order). Compute size, alignment, field offsets.
- **Enum:** compute layout: tag + union of variant payloads.
- **Interface:** no layout per se; record method set.
- **Const `const N = expr`:** evaluate `expr` at compile time. Requires a small compile-time interpreter. The result is a `ConstValue`.
- **Global `var x: T = expr`:** resolve `T`; do NOT evaluate `expr` at compile time (it's a runtime initializer). Record the type and mark the global as initialized by the init function.
- **Function `fn f(...) -> ...`:** resolve parameter and return types. Do NOT typecheck the body yet (that's a separate pass).

After this pass, all types have layouts, all constants have values, and all function signatures are resolved.

**Then** the typechecker runs on function bodies and constant expressions that reference constants.

Wait — constant expressions already need to be evaluated in the topological pass, because types depend on them. So you need the compile-time evaluator *before* the typechecker completes. Hmm.

Let me restructure.

---

## 4. The Actual Pass Order

Here's the full pipeline for a package:

```
1. Parse all files in the package → untyped ASTs
2. Collect top-level declarations → package scope
3. Resolve all top-level names in signatures and constant expressions
4. Build dependency graph
5. Topological sort → declaration order
6. For each declaration in order:
     a. Resolve its internal references (types, constants)
     b. If const: evaluate its expression to a ConstValue
     c. If type: compute layout
     d. If function: record signature
7. Typecheck function bodies and global initializers
8. Lower to IR
```

Step 6b requires a compile-time evaluator. Step 6c requires knowing types of fields, which requires type resolution (step 6a) and constant evaluation (6b). Step 6d just needs signatures.

**The compile-time evaluator is the hard part.** Let's look at it.

---

## 5. The Compile-Time Evaluator

You need to evaluate constant expressions at compile time. These expressions can reference:

- Integer, float, bool, string literals.
- Other constants.
- Arithmetic, comparison, logical operators.
- `sizeof(T)`, `alignof(T)`, `offsetof(T, field)` builtins.
- Struct literals (for `EMPTY = Buffer{...}`).
- Array literals.
- Function calls **if and only if** the function is a `comptime` function (v2) — for v1, no function calls in constants.

**Restriction for v1:** constant expressions are limited to:
- Literals
- References to other constants
- Arithmetic/bitwise/comparison/logical operators
- The builtins `sizeof`, `alignof`, `offsetof`
- Type conversions between numeric types
- Array/struct literals (whose fields are themselves constant)

No function calls. No control flow. No loops. No `if`. This makes the evaluator a simple recursive tree walker:

```rust
fn eval_const(&self, expr: &TypedExpr, env: &ConstEnv) -> ConstValue {
    match &expr.kind {
        IntLit(n) => ConstValue::Int(*n, expr.ty),
        FloatLit(f) => ConstValue::Float(*f, expr.ty),
        BoolLit(b) => ConstValue::Bool(*b),
        StringLit(s) => ConstValue::String(s.clone()),
        ConstRef(id) => env.get(id),  // already evaluated
        Binary(op, a, b) => {
            let av = self.eval_const(a, env);
            let bv = self.eval_const(b, env);
            apply_binop(*op, av, bv)
        }
        Builtin(BuiltinKind::SizeOf, ty) => ConstValue::Int(self.layout_of(*ty).size as i64, U64),
        Builtin(BuiltinKind::AlignOf, ty) => ConstValue::Int(self.layout_of(*ty).align as i64, U64),
        Builtin(BuiltinKind::OffsetOf, ty, field) => ConstValue::Int(self.offset_of(*ty, field) as i64, U64),
        ArrayLit(elems) => ConstValue::Array(elems.iter().map(|e| self.eval_const(e, env)).collect()),
        StructLit(fields) => ConstValue::Struct(fields.iter().map(|(n, e)| (n.clone(), self.eval_const(e, env))).collect()),
        _ => error("not a constant expression"),
    }
}
```

**Typechecking constant expressions:** You must typecheck them before evaluating. This requires the typechecker to be available during the topological pass. Circular dependency again.

**Resolution:** Split the typechecker into two phases:

- **`check_const_expr(expr, expected_ty) -> TypeId`:** a restricted typechecker for constant expressions. Only allows the constructs above. Does not need to know about function bodies, control flow, or local variables.
- **`check_function_body(fn) -> ()`:** the full typechecker. Uses `check_const_expr` for constant subexpressions.

The restricted typechecker is a subset of the full one. You can implement it as a flag on the full typechecker: `in_const_context: bool`. When true, reject constructs that aren't allowed in constants.

This is clean. Do it.

---

## 6. Handling the Builtins

`sizeof`, `alignof`, `offsetof` need layout information. But layout information depends on constant evaluation (for array sizes). And constant evaluation needs layout information (for these builtins). Circular.

**Resolution:** The dependency graph handles it. `sizeof(T)` in a constant expression creates a dependency edge from the constant to `T`. So `T` is laid out before the constant is evaluated. The topological order ensures this.

But what if `T`'s layout depends on the constant? Then there's a genuine cycle, and the topological sort reports it. Correct behavior.

**Builtins must appear in the dependency graph.** When collecting references from a constant expression, `sizeof(T)` contributes a reference to `T`. Same for `alignof` and `offsetof`.

---

## 7. Example Walkthrough

Let's trace the example from §1:

```lang
const N = 16
const M = N * 2

type Buffer = struct {
    data: [M]u8
    len:  usize
}

const EMPTY = Buffer{ data: [M]u8{}, len: 0 }

type Aligned = struct {
    bytes: [alignof(Buffer)]u8
}
```

**Phase A: Collect names.**
Package scope: `N`, `M`, `Buffer`, `EMPTY`, `Aligned`.

**Phase B: Resolve references and build graph.**
- `N` → no deps.
- `M` → depends on `N`.
- `Buffer` → depends on `M` (via `[M]u8`).
- `EMPTY` → depends on `Buffer`, `M`.
- `Aligned` → depends on `Buffer` (via `alignof(Buffer)`).

Graph:
```
N ← M ← Buffer ← EMPTY
              ↖
              Aligned
```

**Topological sort:** `N, M, Buffer, {EMPTY, Aligned}` (the last two are independent).

**Evaluation:**
1. `N = 16`.
2. `M = 16 * 2 = 32`.
3. `Buffer`: layout = `{ data: [32]u8, len: usize }`. Say `usize` is 8 bytes: size = 32 + 8 = 40, align = 8.
4. `EMPTY`: evaluate `Buffer{ data: [32]u8{}, len: 0 }` → a `ConstValue::Struct` with the fields. (Note: `[32]u8{}` is a zero-initialized array; support this as sugar for "all elements default".)
5. `Aligned`: layout = `{ bytes: [alignof(Buffer)]u8 }` = `{ bytes: [8]u8 }`. Size = 8, align = 1.

Correct.

---

## 8. Cross-Package Dependencies

Packages complicate things. Package A imports package B. `B`'s exported declarations must be hoisted before `A`'s.

**Rule:** Build the dependency graph at the *program* level, not the package level. Nodes are declarations; edges can cross packages. Topological sort globally. Cycle detection catches cross-package cycles too.

**But** cross-package cycles are usually a *good* thing to forbid (Go does). If `A` imports `B` and `B` imports `A`, report an error immediately during import resolution — you don't even need to build the declaration graph.

So the algorithm is:

1. Resolve imports. Build package-level dependency graph. Topological sort packages. Detect cycles → error.
2. For each package in topological order, hoist its declarations.
3. Cross-package references only go from later packages to earlier packages, so they're always resolvable.

**This is much simpler than a global graph.** Do this.

---

## 9. Incremental Compilation Considerations

If you want incremental compilation (recompile only what changed), you need:

- A content hash of each declaration (its AST + the hashes of its dependencies).
- A cache mapping hash → compiled artifact (layout, const value, IR).

When a declaration changes, recompute its hash. If the hash changed, invalidate its dependents. Recompute in topological order.

The dependency graph is exactly what you need. Store it persistently (or recompute it — it's cheap).

**For v1, don't do incremental.** Compile the whole package each time. Add incremental when compile times become a problem.

---

## 10. Where This Plugs into the Earlier Plan

From the earlier typechecking/lowering plan:

- **Resolver (Pass 1)** becomes:
  - Pass 1a: Collect top-level names.
  - Pass 1b: Build dependency graph.
  - Pass 1c: Topological sort.
  - Pass 1d: For each declaration in order, resolve references, evaluate constants, compute layouts.
  - Pass 1e: Resolve and typecheck function bodies (can be parallelized).

- **Typechecker (Pass 2)** runs after Pass 1e. It can assume:
  - All types have layouts.
  - All constants have values.
  - All signatures are resolved.

- **Lowering (Pass 3)** runs after Pass 2. It can assume everything is typed.

The **compile-time evaluator** is a new component that runs inside Pass 1d. It's used by:
- Constant declarations.
- Array lengths.
- `sizeof`/`alignof`/`offsetof`.
- Enum discriminant values (if you have them).
- Global initializers (partially — only if the initializer is a constant expression).

---

## 11. Subtleties and Edge Cases

### 11.1 Forward References in Function Bodies

Function bodies can reference declarations that come *later* in source order, as long as there's no cycle. The topological sort handles this naturally: by the time you typecheck function bodies (Pass 1e), all declarations are hoisted. So:

```lang
fn f() -> i32 { return g() }
fn g() -> i32 { return 42 }
```

Works. `f`'s body is typechecked after `g` is resolved.

### 11.2 Recursive Types

```lang
type Node = struct {
    value: i32
    next:  ?*Node
}
```

`Node` depends on `Node` (via `?*Node`). This is a self-loop in the graph. **It's not an error**, because pointer types don't require the pointee's layout to be known — a pointer is always pointer-sized.

**Rule:** When collecting dependencies, **do not create an edge for a pointer type's pointee.** A pointer to `T` doesn't require `T` to be laid out. This allows recursive types via pointers.

Same for slices: `[]T` is `{ptr, len}`, doesn't require `T`'s layout. Same for function pointers.

But `[N]T` *does* require `T`'s layout (you need to know the element size to compute the array size). And struct fields require their types' layouts. And enum payloads require their types' layouts.

So:

- Pointer/slice/function-pointer: no edge to pointee.
- Array: edge to element type and to length expression.
- Struct field: edge to field type.
- Enum variant field: edge to field type.

This makes recursive types work naturally.

### 11.3 Mutually Recursive Types via Pointers

```lang
type A = struct { b: ?*B }
type B = struct { a: ?*A }
```

No edges (pointers), so no cycle. Both types get laid out independently. Correct.

### 11.4 Constants in Type Positions

```lang
type Arr = [N]u8
```

`Arr` depends on `N`. `N` must be a constant. The dependency graph captures this.

What if `N` is not a constant?

```lang
var N = 16
type Arr = [N]u8    // error
```

The dependency graph has an edge from `Arr` to `N`, but `N` is a global variable, not a constant. During hoisting, when you try to evaluate `Arr`'s layout, you need `N`'s value, but globals don't have compile-time values. Error: "array length must be a constant."

Detect this in the resolver: when collecting references from a type declaration, if a referenced declaration is a global variable (not a constant), report an error.

### 11.5 Function Calls in Constants

You banned them for v1. But when you add `comptime`, you need to:

- Evaluate the function at compile time.
- Detect infinite recursion (depth limit).
- Cache results.

The dependency graph must include edges from the constant to the `comptime` function and from the function to any constants it references. This gets hairy. **Defer to v2.** For v1, no function calls in constants.

### 11.6 Order of Error Reporting

When there are cycles, report the cycle, not the downstream errors. Example:

```lang
const A = B + 1
const B = A + 1
const C = A + 1
```

Report "cycle: A → B → A". Don't report "C references undefined A" — `A` is defined, it's just in a cycle. Suppress errors for declarations in the cycle.

**Implementation:** when a cycle is detected, mark all declarations in the cycle as `ERROR_TYPE`. Then continue hoisting. Errors that depend on cycle members will see `ERROR_TYPE` and not produce spurious messages.

### 11.7 Duplicate Declarations

```lang
const N = 16
const N = 32    // error
```

Detect during Phase A (name collection). Two declarations with the same name in the same package scope → error.

### 11.8 Shadowing Builtins

```lang
const len = 16    // shadows builtin len?
```

Decide: builtins are keywords or predeclared identifiers? I recommend **predeclared identifiers** (like Go). Users can shadow them. The package scope has `len`, `sizeof`, etc.; a user declaration with the same name shadows them.

But then `type T = [len]u8` might refer to the user's `len`, not the builtin. This is fine — it's consistent with lexical scoping. Document it.

---

## 12. Implementation Checklist

1. **Extend the AST** with a `DeclId` field on each top-level declaration. Assign IDs during Phase A.

2. **Write `collect_top_level_names`.** Walk all package files, insert each top-level name into the package scope. Detect duplicates.

3. **Write `collect_refs(decl) -> Vec<DeclId>`.** Walk a declaration's AST, resolve identifiers, collect referenced declaration IDs. Handle pointer/slice/array distinctions correctly (no edge for pointee of pointer).

4. **Build the dependency graph.** Nodes = declarations. Edges = references.

5. **Topological sort with cycle detection.** Use DFS with a state array and a stack for cycle extraction.

6. **Write the compile-time evaluator.** Tree walker over typed constant expressions. Handles literals, references to already-evaluated constants, arithmetic, comparisons, builtins.

7. **Write the layout calculator.** Given a struct/enum/array type, compute size, alignment, field offsets. Requires the layouts of field types (already computed).

8. **The hoisting loop:** iterate over declarations in topological order. Dispatch by kind:
   - Type: compute layout, store in type table.
   - Const: typecheck as constant expression, evaluate, store value.
   - Func: resolve signature, store.
   - Global: resolve type, store.

9. **Detect invalid dependencies:** a type that depends on a global variable; a constant that depends on a global variable; a constant expression that calls a function (v1 restriction).

10. **After hoisting, typecheck function bodies.** Use the full typechecker. Constant subexpressions use the restricted typechecker.

---

## 13. A Note on Simplicity

You might be tempted to avoid the topological sort by requiring source order (like C). **Don't.** Source order is a wart that users hate, and it doesn't even solve the problem (constants can reference later constants via forward declarations). The topological sort is ~100 lines of code and eliminates a whole class of user frustration.

Similarly, you might be tempted to allow only "constants before types" or some other partial order. **Don't.** A full topological sort is no harder than a partial one and handles all cases uniformly.

The one thing you *should* restrict is **function calls in constants**. This is where the complexity explodes. For v1, no function calls in constants. Constants are expressions over literals, other constants, operators, and builtins. That's it.

---

## 14. Summary

The algorithm is:

1. **Collect all top-level names** in one pass. Detect duplicates.
2. **Build a dependency graph** over declarations by resolving references. Pointers don't create dependencies on their pointees.
3. **Topologically sort.** Detect cycles. Report them precisely.
4. **Process declarations in order:**
   - Constants: typecheck as constant expressions, evaluate.
   - Types: compute layout (requires field types' layouts).
   - Functions: resolve signatures only.
   - Globals: resolve types only.
5. **Then typecheck function bodies.**

The compile-time evaluator handles constants. The layout calculator handles types. The dependency graph handles ordering. Cycles are errors. Recursive types via pointers are allowed.

This is the standard approach. It's what Go, Swift, Rust, and Zig all do (with variations). It's implementable in a few hundred lines. Do it.
