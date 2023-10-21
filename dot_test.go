package dot

import (
	"fmt"
	"go/format"
	"reflect"
	"testing"
)

func TestStruct_ExecuteTemplateString(t *testing.T) {
	st := Struct{
		Name:       "TestStruct",
		TypeParams: "[T any, R string]",
		Fields: []Field{
			{
				Name:     "Name",
				Type:     NewType("*fmt.Stringer"),
				Tag:      `json:"name"`,
				Comments: []string{"the name"},
			},
			{
				Name:     "Age",
				Type:     Int,
				Tag:      `json:"age"`,
				Comments: []string{"the age"},
			},
		},
		Directives: []string{"go:generate go fmt ./..."},
		Comments:   []string{"this is to test the Struct"},
		Templates:  []string{TemplateStruct},
		Variables:  nil,
	}
	text, err := st.ExecuteTemplateString(st.Templates[0])
	if err != nil {
		t.Error(err)
		return
	}
	b, err := format.Source([]byte(text))
	if err != nil {
		t.Errorf("error format source code: %v, source: %s", err, text)
		return
	} else {
		fmt.Println(string(b))
	}
}

func TestNewType(t *testing.T) {
	type args struct {
		fullTypeName string
	}
	tests := []struct {
		name string
		args args
		want Type
	}{
		{
			name: "single",
			args: args{
				fullTypeName: "string",
			},
			want: TypeImpl{Name: "string"},
		},
		{
			name: "with package",
			args: args{
				fullTypeName: "fmt.Stringer",
			},
			want: TypeImpl{PackageName: "fmt", Name: "Stringer"},
		},
		{
			name: "with pointer",
			args: args{
				fullTypeName: "*string",
			},
			want: TypeImpl{Stars: 1, Name: "string"},
		},
		{
			name: "with pointer and package",
			args: args{
				fullTypeName: "*fmt.Stringer",
			},
			want: TypeImpl{Stars: 1, PackageName: "fmt", Name: "Stringer"},
		},
		{
			name: "with package and many pointers",
			args: args{
				fullTypeName: "**fmt.Stringer",
			},
			want: TypeImpl{Stars: 2, PackageName: "fmt", Name: "Stringer"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewType(tt.args.fullTypeName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewType() = %v, want %v", got, tt.want)
			} else {
				fmt.Println(got)
			}
		})
	}
}
