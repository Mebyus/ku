package char

import (
	"testing"
)

func TestUnescape(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
		ok   bool
	}{
		{
			name: "1 empty string",
			s:    "",
			want: "",
			ok:   true,
		},
		{
			name: "2 one char no escapes",
			s:    "a",
			want: "a",
			ok:   true,
		},
		{
			name: "3 no escapes",
			s:    "abc 123",
			want: "abc 123",
			ok:   true,
		},
		{
			name: "4 line feed",
			s:    `\n`,
			want: "\n",
			ok:   true,
		},
		{
			name: "5 backslash",
			s:    `\\`,
			want: "\\",
			ok:   true,
		},
		{
			name: "6 carriage return after string",
			s:    `abc\r`,
			want: "abc\r",
			ok:   true,
		},
		{
			name: "7 bad escape",
			s:    `\`,
			want: "",
			ok:   false,
		},
		{
			name: "8 bad escape",
			s:    `abc \h 123`,
			want: "abc  123",
			ok:   false,
		},
		{
			name: "9 null escape",
			s:    `\0`,
			want: "\x00",
			ok:   true,
		},
		{
			name: "10 hex escape",
			s:    `\xA4`,
			want: "\xA4",
			ok:   true,
		},
		{
			name: "11 bad hex escape",
			s:    `\xA`,
			want: "",
			ok:   false,
		},
		{
			name: "12 hex escape",
			s:    `abc \x00123`,
			want: "abc \x00123",
			ok:   true,
		},
		{
			name: "13 bad hex escape",
			s:    `abc\xAG123`,
			want: "abc123",
			ok:   false,
		},
		{
			name: "14 apostrophe",
			s:    `\'A\'`,
			want: "'A'",
			ok:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := Unescape(tt.s)
			if ok != tt.ok || got != tt.want {
				t.Errorf("Unescape() = (\"%s\", %v), want (\"%s\", %v)", got, ok, tt.want, tt.ok)
			}
		})
	}
}
