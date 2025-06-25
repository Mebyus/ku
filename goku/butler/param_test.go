package butler

import (
	"fmt"
	"reflect"
	"testing"
)

type testBox1 struct{}

var _ ParamBox = &testBox1{}

func (b *testBox1) Get(name string) *Param {
	panic("no params expected")
}

func (b *testBox1) Apply(p *Param) error {
	panic("no params expected")
}

func (b *testBox1) Params() []Param {
	return nil
}

type testBox2 struct {
	a string
	b bool
}

var _ ParamBox = &testBox2{}

func (b *testBox2) Get(name string) *Param {
	switch name {
	case "a":
		return &Param{
			Name:    "a",
			Default: "",
			Kind:    String,
			val:     b.a,
		}
	case "b":
		return &Param{
			Name:    "b",
			Default: false,
			Kind:    Boolean,
			val:     b.a,
		}
	default:
		panic(fmt.Sprintf("unexpected param \"%s\"", name))
	}
}

func (b *testBox2) Apply(p *Param) error {
	switch p.Name {
	case "a":
		b.a = p.Str()
	case "b":
		b.b = p.Bool()
	default:
		panic(fmt.Sprintf("unexpected param \"%s\"", p.Name))
	}
	return nil
}

func (b *testBox2) Params() []Param {
	return []Param{
		{
			Name:    "a",
			Default: "",
			Kind:    String,
		},
		{
			Name:    "b",
			Default: false,
			Kind:    Boolean,
		},
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		box     ParamBox
		args    []string
		want    []string
		wantBox ParamBox
		wantErr bool
	}{
		{
			name:    "1 nil box no args",
			box:     nil,
			args:    nil,
			want:    nil,
			wantBox: nil,
		},
		{
			name:    "2 nil box",
			box:     nil,
			args:    []string{"a", "b"},
			want:    []string{"a", "b"},
			wantBox: nil,
		},
		{
			name:    "3 no params box no args",
			box:     &testBox1{},
			args:    nil,
			want:    nil,
			wantBox: &testBox1{},
		},
		{
			name:    "4 no params box",
			box:     &testBox1{},
			args:    []string{"a", "b"},
			want:    []string{"a", "b"},
			wantBox: &testBox1{},
		},
		{
			name:    "5 string param no args",
			box:     &testBox2{},
			args:    nil,
			want:    nil,
			wantBox: &testBox2{},
		},
		{
			name:    "6 string param no matching args",
			box:     &testBox2{},
			args:    []string{"a", "b"},
			want:    []string{"a", "b"},
			wantBox: &testBox2{},
		},
		{
			name:    "6 string param matching args",
			box:     &testBox2{},
			args:    []string{"--a=hello", "b"},
			want:    []string{"b"},
			wantBox: &testBox2{a: "hello"},
		},
		{
			name:    "7 string param matching args",
			box:     &testBox2{},
			args:    []string{"--a", "hello", "b"},
			want:    []string{"b"},
			wantBox: &testBox2{a: "hello"},
		},
		{
			name:    "8 boolean param matching args",
			box:     &testBox2{},
			args:    []string{"--b"},
			want:    nil,
			wantBox: &testBox2{b: true},
		},
		{
			name:    "9 boolean param matching args",
			box:     &testBox2{},
			args:    []string{"--b=true"},
			want:    nil,
			wantBox: &testBox2{b: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.box, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}

			box := tt.box
			if !reflect.DeepEqual(box, tt.wantBox) {
				t.Errorf("ParamBox = %v, want %v", box, tt.wantBox)
			}
		})
	}
}
