package cc

var warningFlags = []string{
	"-Wall",
	"-Wextra",
	"-Wconversion",
	"-Wunreachable-code",
	"-Wshadow",
	"-Wundef",
	"-Wfloat-equal",
	"-Wformat=0",
	"-Wpointer-arith",
	"-Winit-self",
	"-Wuninitialized",
	"-Wduplicated-branches",
	"-Wduplicated-cond",
	"-Wdouble-promotion",
	// "-Wnull-dereference",
	"-Wstrict-prototypes",
	// "-Wvla",
	"-Wpointer-sign",
	"-Wswitch-default",
	"-Wshadow=local",

	// "-Wno-main",
	"-Wno-shadow",
	"-Wno-unused-parameter",
	"-Wno-unused-function",
	"-Wno-unused-const-variable",
}

var otherFlags = []string{
	"-Werror",
	"-pipe",
	// "-fanalyzer",
}

var codegenFlags = []string{
	"-fwrapv",
	"-funsigned-char",
	"-fno-asynchronous-unwind-tables",
	"-fomit-frame-pointer",
	"-fno-stack-protector",
}

const debugCompilerOptimizations = "g"
const testCompilerOptimizations = "1"
const safeCompilerOptimizations = "2"
const fastCompilerOptimizations = "fast"

const debugInfoFlag = "-ggdb"

const maxErrorsFlag = "-fmax-errors=1"
