package token

import "testing"

func TestLiteral(t *testing.T) {
	for k := EOF + 1; k < _keyend; k += 1 {
		if k == _keystart {
			continue
		}

		if ktext[k] == "" {
			t.Errorf("token kind (=%d) has empty text literal", k)
		}
	}
}
