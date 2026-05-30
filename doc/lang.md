# Ku programming language

## Source code

Ku source code (sometimes also referred as "source text" or just "text") comes in plain text files
typically with ".ku" extension. Compiler assumes that each file contains utf-8 text, although it
only uses ascii for language syntax and takes "hands off" approach when it reads string literals
or comments, treating them as opaque byte sequences. Thus language only accepts identifiers which
start with latin letter or underscore and contains only latin letters, underscores and decimal
digits.

### Units

Language, compiler and build system primarily works with collection of files called "unit" as
smallest unit of compilation (however it is possible to dissect individual files via low-level
commands), all files within the same unit are processed and analyzed together at once. Build system
in particular treats all "*.ku" files in the same directory as one unit.

Top-level symbols defined in unit are shared between all unit files. Thus each unit creates a separate
namespace. Units may access symbols from other units via imports.

### Imports

A unit can import other units via top-level import construct:

```ku
// Import unit "fmt" from standard library and bind it
// to this unit's local top-level symbol "fmt"
import std (
    fmt -> "fmt"
)

// Import local unit "zzz/a" and bind it to this unit's local top-level symbol "a"
import (
    a -> "zzz/a"
)
```

When importing a unit it is specified by "import string" and "import origin".
How exactly import string and origin is translated to actual imported unit source
code is determined by build system. In general import origin should specify a common
place for unit lookup. Native build system recognizes 3 distinct import origins:

- std - standard library
- pkg - third-party libraries
- loc - (or just empty origin) local origin, meaning that units are located in current package

Native build system compiles and builds a whole program starting from a specified unit.
It then locates all imported units, forming a dependecy graph between units.
Correct program cannot have import cycles.

## Implementation details

Although compiler can be used to digest any file and treat it as Ku source code, the native build system
is more restrictive, it only recognizes files with names in form "*.ku" (and part before ".ku" extension
must not be empty). Build system also skips empty or too large (above 64 MB in size) source files.

Compiler has maximum byte length limit for most tokens (4096 bytes) including identifiers and excluding
raw string literals and comments.
