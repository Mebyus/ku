package token

import "fmt"

// Compare returns a non-nil error if given tokens are not equal.
// Comparison ignores token pins.
func Compare(a, b Token) error {
	if a.Kind != b.Kind {
		return fmt.Errorf("kind: got %s, want %s", a.Kind, b.Kind)
	}
	if a.Val != b.Val {
		return fmt.Errorf("val: got %d, want %d", a.Val, b.Val)
	}
	if a.Data != b.Data {
		return fmt.Errorf("data: got %s, want %s", a.Data, b.Data)
	}
	return nil
}
