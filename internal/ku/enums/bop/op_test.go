package bop

import (
	"testing"
)

func TestKind(t *testing.T) {
	// Test checks that precedence and power of all valid operators are >= 1.

	for i := range num {
		k := Kind(i)
		pre := k.Precedence()
		p := k.Power()

		if pre < 1 {
			t.Errorf("Precedence(%s) = %d", k, pre)
		}
		if p < 1 {
			t.Errorf("Power(%s) = %d", k, p)
		}
	}
}
