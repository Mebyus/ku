package kvx

import (
	"bytes"
	"testing"
)

func TestEncode(t *testing.T) {
	tests := []struct {
		name string

		text string
		data string

		gsize uint32
	}{
		{
			name: "1 empty program",

			text:  "",
			data:  "",
			gsize: 0,
		},
		{
			name: "2 only text",

			text:  "abc",
			data:  "",
			gsize: 0,
		},
		{
			name: "3 only data",

			text:  "",
			data:  "123",
			gsize: 0,
		},
		{
			name: "4 only global",

			text:  "",
			data:  "",
			gsize: 16,
		},
		{
			name: "5 generic",

			text:  "Hello Text Hello Text",
			data:  "Hello Data Hello Data",
			gsize: 9230,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			progIn := Program{
				Text:       []byte(tt.text),
				Data:       []byte(tt.data),
				GlobalSize: tt.gsize,
			}

			var out bytes.Buffer
			err := Encode(&out, &progIn)
			if err != nil {
				t.Errorf("Encode() error = %v", err)
				return
			}

			progOut, err := Decode(bytes.NewReader(out.Bytes()))
			if err != nil {
				t.Errorf("Decode() error = %v", err)
				return
			}

			if !bytes.Equal(progOut.Text, progIn.Text) {
				t.Errorf("Decode() Text = \"%s\", want \"%s\"", progOut.Text, progIn.Text)
			}
			if !bytes.Equal(progOut.Data, progIn.Data) {
				t.Errorf("Decode() Data = \"%s\", want \"%s\"", progOut.Data, progIn.Data)
			}
			if progOut.GlobalSize != progIn.GlobalSize {
				t.Errorf("Decode() GlobalSize = %d, want %d", progOut.GlobalSize, progIn.GlobalSize)
			}
		})
	}
}
