package char

import "testing"

func TestSnakeCase(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{
			name: "1 empty string",
			s:    "",
			want: "",
		},
		{
			name: "2 single letter",
			s:    "a",
			want: "a",
		},
		{
			name: "3 single capital letter",
			s:    "A",
			want: "a",
		},
		{
			name: "4 simple word",
			s:    "aa",
			want: "aa",
		},
		{
			name: "5 simple capital word",
			s:    "Aa",
			want: "aa",
		},
		{
			name: "6 simple capital word",
			s:    "Abc",
			want: "abc",
		},
		{
			name: "7 two words",
			s:    "AbcAbc",
			want: "abc_abc",
		},
		{
			name: "8 two words",
			s:    "Abc_Abc",
			want: "abc_abc",
		},
		{
			name: "9 two words",
			s:    "abc_abc",
			want: "abc_abc",
		},
		{
			name: "10 three words",
			s:    "abcAbcZ",
			want: "abc_abc_z",
		},
		{
			name: "11 two words",
			s:    "ZabcAbc",
			want: "zabc_abc",
		},
		{
			name: "12 three words",
			s:    "zAbcAbc",
			want: "z_abc_abc",
		},
		{
			name: "13 two words",
			s:    "abc_Abc",
			want: "abc_abc",
		},
		{
			name: "14 capital word",
			s:    "ABC",
			want: "abc",
		},
		{
			name: "15 two capital words",
			s:    "ABC_ZABC",
			want: "abc_zabc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SnakeCase(tt.s); got != tt.want {
				t.Errorf("SnakeCase() = %v, want %v", got, tt.want)
			}
		})
	}
}
