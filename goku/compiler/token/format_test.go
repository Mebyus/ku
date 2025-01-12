package token

import "testing"

func TestKind_String(t *testing.T) {
	if len(literal) != int(maxKind) {
		t.Errorf("token kind literals are missing: got %d, want %d", len(literal), maxKind)
	}

	for k := range staticLiteralEnd {
		if literal[k] == "" {
			t.Errorf("token kind (=%d) has empty literal", k)
		}
	}

	// staticLiteralEnd and maxKind are sentinel values, they do not have literals
	k := staticLiteralEnd + 1
	for k < maxKind {
		if literal[k] == "" {
			t.Errorf("token kind (=%d) has empty literal", k)
		}
		k += 1
	}
}
