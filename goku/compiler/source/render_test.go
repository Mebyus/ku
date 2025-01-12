package source

import (
	"bytes"
	"strings"
	"testing"
)

const textInput1 = `
fun add(a: s32, b: s32) => s32 {
	ret a + b;
}
`

const renderOutput1 = `
* >>>> example.ku:2:6
  |
2 |	fun add(a: s32, b: s32) => s32 {
3 |		ret a + b;
* >>        ^^^^^
4 |	}
5 |
`

func TestRender(t *testing.T) {
	text1 := NewText("example.ku", []byte(textInput1))

	tests := []struct {
		name  string
		opts  RenderParams
		wantW string
	}{
		{
			name: "1",
			opts: RenderParams{
				Text: text1,

				Window: WindowParams{
					MaxLinesBefore: 2,
					MaxLinesAfter:  2,
				},
			},
			wantW: strings.TrimSpace(renderOutput1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			err := Render(w, tt.opts)
			if err != nil {
				t.Errorf("RenderSpan() error = %v", err)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("RenderSpan()\n\n%v\n\nwant\n\n%v", gotW, tt.wantW)
			}
		})
	}
}

func TestFindOffset(t *testing.T) {
	tests := []struct {
		name   string
		text   string
		pos    string
		want   uint32
		wantOk bool
	}{
		{
			name:   "1",
			text:   textInput1,
			pos:    "1:1",
			want:   0,
			wantOk: true,
		},
		{
			name:   "2",
			text:   textInput1,
			pos:    "1:2",
			want:   0,
			wantOk: false,
		},
		{
			name:   "3",
			text:   textInput1,
			pos:    "2:28",
			want:   28,
			wantOk: true,
		},
		{
			name:   "4",
			text:   textInput1,
			pos:    "2:33",
			want:   33,
			wantOk: true,
		},
		{
			name:   "5",
			text:   textInput1,
			pos:    "2:1",
			want:   1,
			wantOk: true,
		},
		{
			name:   "6",
			text:   textInput1,
			pos:    "4:1",
			want:   46,
			wantOk: true,
		},
		{
			name:   "7",
			text:   textInput1,
			pos:    "4:2",
			want:   47,
			wantOk: true,
		},
		{
			name:   "8",
			text:   textInput1,
			pos:    "2:34",
			want:   0,
			wantOk: false,
		},
		{
			name:   "9",
			text:   textInput1,
			pos:    "5:1",
			want:   48,
			wantOk: true,
		},
		{
			name:   "10",
			text:   textInput1,
			pos:    "5:3",
			want:   0,
			wantOk: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			text := []byte(tt.text)
			got, gotOk := FindOffset(text, MustParseTextPos(tt.pos))
			if got != tt.want || gotOk != tt.wantOk {
				t.Errorf("FindOffset() got = (%v, %v), want (%v, %v)", got, gotOk, tt.want, tt.wantOk)
			}
			if tt.wantOk {
				gotPos := FindTextPos(text, tt.want).String()
				if gotPos != tt.pos {
					t.Errorf("FindTextPos() got = \"%s\", want \"%s\"", gotPos, tt.pos)
				}
			}
		})
	}
}

func TestFindLineOffset(t *testing.T) {
	tests := []struct {
		name   string
		text   string
		offset uint32
		want   uint32
	}{
		{
			name:   "1",
			text:   "",
			offset: 0,
			want:   0,
		},
		{
			name:   "2",
			text:   " ",
			offset: 0,
			want:   0,
		},
		{
			name:   "3",
			text:   " ",
			offset: 1,
			want:   0,
		},
		{
			name:   "4",
			text:   "\n",
			offset: 0,
			want:   0,
		},
		{
			name:   "5",
			text:   "\n",
			offset: 1,
			want:   1,
		},
		{
			name:   "6",
			text:   "abc\n",
			offset: 1,
			want:   0,
		},
		{
			name:   "7",
			text:   "abc\n",
			offset: 3,
			want:   0,
		},
		{
			name:   "8",
			text:   "abc\n",
			offset: 4,
			want:   4,
		},
		{
			name:   "9",
			text:   "abc\nabc",
			offset: 4,
			want:   4,
		},
		{
			name:   "10",
			text:   "abc\nabc",
			offset: 7,
			want:   4,
		},
		{
			name:   "11",
			text:   "abc\nabc\n",
			offset: 7,
			want:   4,
		},
		{
			name:   "12",
			text:   "abc\nabc\n",
			offset: 8,
			want:   8,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FindLineOffset([]byte(tt.text), tt.offset); got != tt.want {
				t.Errorf("FindLineOffset() = %v, want %v", got, tt.want)
			}
		})
	}
}
