package token

import (
	"testing"

	"github.com/mebyus/ku/goku/compiler/baselex"
)

func TestFromBaseKind(t *testing.T) {
	tests := []struct {
		base uint32
		want Kind
	}{
		{
			base: baselex.Illegal,
			want: Illegal,
		},
		{
			base: baselex.Word,
			want: Word,
		},
		{
			base: baselex.BinInteger,
			want: BinInteger,
		},
		{
			base: baselex.OctInteger,
			want: OctInteger,
		},
		{
			base: baselex.DecInteger,
			want: DecInteger,
		},
		{
			base: baselex.HexInteger,
			want: HexInteger,
		},
		{
			base: baselex.DecFloat,
			want: DecFloat,
		},
		{
			base: baselex.Rune,
			want: Rune,
		},
		{
			base: baselex.Rune,
			want: Rune,
		},
	}
	for _, tt := range tests {
		t.Run(tt.want.String(), func(t *testing.T) {
			if got := FromBaseKind(tt.base); got != tt.want {
				t.Errorf("FromBaseKind() = %v, want %v", got, tt.want)
			}
		})
	}
}
