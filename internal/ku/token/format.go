package token

import (
	"fmt"
)

// String returns token textual representation (without position info).
//
// Should be used mainly for error reporting.
func (t *Token) String() string {
	cat := t.Kind.Category()
	switch t.Kind {
	case INV, EOF:
		return cat
	case Word:
		return fmt.Sprintf("%s %s", cat, t.Data)
	case Integer:
		if t.Data != "" {
			return fmt.Sprintf("%s %s", cat, t.Data)
		}

		switch t.Val {
		case DecInt:
			return fmt.Sprintf("%s %d", cat, t.Val)
		case BinInt:
			return fmt.Sprintf("%s %b", cat, t.Val)
		case OctInt:
			return fmt.Sprintf("%s %o", cat, t.Val)
		case HexInt:
			return fmt.Sprintf("%s %X", cat, t.Val)
		default:
			return fmt.Sprintf("%s %d", cat, t.Val)
		}
	case String:
		return fmt.Sprintf("%s \"%s\"", cat, t.Data) // TODO: need to escape this
	}

	return cat // TODO: literals for keywords and other things
}
