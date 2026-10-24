package ir

type Fun struct {
	// Should always contain return as last statement.
	// Thus it should contain at least one statement.
	Body []Stm

	Nloc uint32 // total number of required local slots (parameters + local variables)
	Narg uint32 // number of required local argument slots
	Nret uint32 // number of returned values (0 or 1 for now)
}
