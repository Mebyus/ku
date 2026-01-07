package sm

import (
	"fmt"
	"testing"
)

func Test_encodeVarUint32(t *testing.T) {
	tests := []struct {
		name string
		v    uint32
	}{
		{
			name: "1",
			v:    0,
		},
		{
			name: "2",
			v:    1,
		},
		{
			name: "3",
			v:    127,
		},
		{
			name: "4",
			v:    128,
		},
		{
			name: "5",
			v:    129,
		},
		{
			name: "6",
			v:    130,
		},
		{
			name: "7",
			v:    0xFFF7A3,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s v=%d", tt.name, tt.v), func(t *testing.T) {
			var buf [16]byte
			v := tt.v
			n := putVarUint32(buf[:], v)
			gotV, gotN := getVarUint32(buf[:])
			if gotV != v || gotN != n {
				t.Errorf("getVarUint32() = (%d, %d); want (%d, %d)", gotV, gotN, v, n)
			}
		})
	}
}
