package sm_test

import (
	"testing"

	"github.com/mebyus/ku/goku/compiler/sm"
)

func TestCheckImportString(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		wantErr bool
	}{
		{
			name:    "1",
			s:       "",
			wantErr: true,
		},
		{
			name:    "2",
			s:       "/",
			wantErr: true,
		},
		{
			name:    "3",
			s:       " ",
			wantErr: true,
		},
		{
			name:    "4",
			s:       "\n",
			wantErr: true,
		},
		{
			name:    "5",
			s:       ".",
			wantErr: true,
		},
		{
			name:    "6",
			s:       "a/",
			wantErr: true,
		},
		{
			name:    "7",
			s:       "/a",
			wantErr: true,
		},
		{
			name:    "8",
			s:       "a/b/",
			wantErr: true,
		},
		{
			name:    "9",
			s:       "/a/b",
			wantErr: true,
		},
		{
			name:    "10",
			s:       "a//b",
			wantErr: true,
		},
		{
			name:    "11",
			s:       "a/ /b",
			wantErr: true,
		},
		{
			name:    "12",
			s:       "a / b",
			wantErr: true,
		},
		{
			name:    "13",
			s:       "a/../b",
			wantErr: true,
		},
		{
			name:    "14",
			s:       "a../b",
			wantErr: true,
		},
		{
			name:    "15",
			s:       "../a/b",
			wantErr: true,
		},
		{
			name:    "16",
			s:       "..a",
			wantErr: true,
		},
		{
			name: "17",
			s:    "a",
		},
		{
			name: "18",
			s:    "abc",
		},
		{
			name: "19",
			s:    "a/bc",
		},
		{
			name: "20",
			s:    "a/b/cd",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := sm.CheckImportString(tt.s)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("CheckImportString() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("CheckImportString() succeeded unexpectedly")
			}
		})
	}
}
