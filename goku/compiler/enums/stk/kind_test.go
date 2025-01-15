package stk

import "testing"

func TestKind(t *testing.T) {
	for k := range maxKind {
		s := k.String()
		if s == "" {
			t.Errorf("statement kind (=%d) has empty literal", k)
		}
	}
}
