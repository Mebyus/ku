package stg

// Hint carries context hint for translating expressions.
type Hint struct {
	// Most close enum type.
	Enum *Enum
}

func (h *Hint) lookupDotName(name string) *EnumEntry {
	if h.Enum == nil {
		return nil
	}
	return h.Enum.m[name]
}
