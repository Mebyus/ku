package ir

type Fun struct {
	// Should always contain return as last statement.
	// Thus it should contain at least one statement.
	Body []Stm

	Narg uint32 // number of required arguments
	Nret uint32 // number of returned values (0 or 1 for now)
}
