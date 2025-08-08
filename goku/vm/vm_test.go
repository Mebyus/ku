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
			name: "1 empty program",
			code: code1,
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

			if exit.Status != tt.want.Status {
				t.Errorf("exit.Status = %d, want %d", exit.Status, tt.want.Status)
			}
		})
	}
}

const code1 = `
#fun start {}
`

const code2 = `
#fun start {
	inc		r0;
}
`

const code3 = `
#fun start {
	halt;
}
`
