package stg

import "github.com/mebyus/ku/goku/compiler/diag"

func (x *TypeIndex) CheckAssign(want *Type, exp Exp) diag.Error {
	typ := exp.Type()
	if typ == want {
		return nil
	}

	return nil
}
