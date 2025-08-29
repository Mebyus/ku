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
				Error: &RuntimeError{Code: ErrorTextEnd},
			},
		},
		{
			name: "2 inc",
			code: code2,
			want: &Exit{
				Error: &RuntimeError{Code: ErrorTextEnd},
			},
		},
		{
			name: "3 halt",
			code: code3,
			want: &Exit{
				Error: nil,
			},
		},
		{
			name: "4 set",
			code: code4,
			want: &Exit{
				Error:  nil,
				Status: 19,
			},
		},
		{
			name: "5 label",
			code: code5,
			want: &Exit{
				Error: nil,
			},
		},
		{
			name: "6 jump",
			code: code6,
			want: &Exit{
				Error: nil,
			},
		},
		{
			name: "7 call",
			code: code7,
			want: &Exit{
				Error:  nil,
				Status: 0x23,
			},
		},
		{
			name: "8 fib",
			code: code8,
			want: &Exit{
				Error:  nil,
				Status: 8,
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
				t.Errorf("exit.Error = (%d) %s", exit.Error.Code, exit.Error)
				return
			}
			if tt.want.Error != nil && exit.Error == nil {
				t.Errorf("exit.Error = <nil>, want (%d) %s", tt.want.Error.Code, tt.want.Error)
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

const code7 = `
#entry start;

#fun start {
	set		#:r0, 0x23;
	call	main;
	halt;
}

#fun main {
	set		#:sc, #:r0;
	ret;
}
`

const code8 = `
#entry start;

#fun start {
	set		#:r0, 6;
	call 	fib;
	set 	#:sc, #:r0;
	halt;
}

// Takes integer argument in #:r0.
// Calculates fibonacci number of that integer.
// Result is returned in #:r0.
// Does not preserve registers.
#fun fib {
	test 		#:r0, 1;
	jump.le		@.exit;

	dec		#:r0;
	push	#:r0;
	call	fib;
	
	pop		#:r1;
	push	#:r0;
	
	set		#:r0, #:r1;
	dec		#:r0;
	call	fib;

	pop		#:r1;
	inc		#:r0, #:r1;

@.exit:
	ret;
}
`
