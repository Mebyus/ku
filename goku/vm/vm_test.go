package vm

import (
	"strings"
	"testing"

	"github.com/mebyus/ku/goku/vm/asm"
)

func TestMachine_Exec(t *testing.T) {
	tests := []struct {
		name string
		code string
		want *Exit
	}{
		{
			name: "1 nop",
			code: code1,
			want: &Exit{
				Error: nil, // TODO: do something about errors comparison
			},
		},
		{
			name: "2 inc",
			code: code2,
			want: &Exit{
				Error: nil, // TODO: do something about errors comparison
			},
		},
		{
			name: "3 halt",
			code: code3,
			want: &Exit{
				Error: nil, // TODO: do something about errors comparison
			},
		},
		{
			name: "4 set",
			code: code4,
			want: &Exit{
				Error:  nil, // TODO: do something about errors comparison
				Status: 19,
			},
		},
		{
			name: "5 label",
			code: code5,
			want: &Exit{
				Error: nil, // TODO: do something about errors comparison
			},
		},
		{
			name: "6 jump",
			code: code6,
			want: &Exit{
				Error: nil, // TODO: do something about errors comparison
			},
		},
	}

	var m Machine
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prog, err := asm.Compile(strings.NewReader(tt.code))
			if err != nil {
				t.Errorf("asm.Compile() error = %v", err)
				return
			}
			exit := m.Exec(prog)

			if tt.want.Error == nil && exit.Error != nil {
				t.Errorf("exit.Error = %v", exit.Error)
				return
			}
			if exit.Status != tt.want.Status {
				t.Errorf("exit.Status = %d, want %d", exit.Status, tt.want.Status)
			}
		})
	}
}

const code1 = `
#entry start;

#fun start {
	nop;
}
`

const code2 = `
#entry start;

#fun start {
	inc		#:r0;
}
`

const code3 = `
#entry start;

#fun start {
	halt;
}
`

const code4 = `
#entry start;

#fun start {
	set		#:sc, 19;
	halt;
}
`

const code5 = `
#entry start;

#fun start {
	nop;

@.label:
	nop;
	halt;
}
`

const code6 = `
#entry start;

#fun start {
	jump	@.label;

@.label:
	nop;
	halt;
}
`
