package exk

import "testing"

func TestKind(t *testing.T) {
	for k := range maxKind {
		s := k.String()
		if s == "" {
			t.Errorf("String(%d) = \"\"", k)
		}
	}
}
