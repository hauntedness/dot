package dot

import (
	"fmt"
	"go/format"
	"strings"
	"testing"
)

type someOption struct {
	name string
	flag bool
}

var WithName = func(name string) func(*someOption) {
	return func(so *someOption) {
		so.name = name
	}
}

var WithFlag = func(flag bool) func(*someOption) {
	return func(so *someOption) {
		so.flag = flag
	}
}

func TestOption(t *testing.T) {
	def := Struct{
		Name: "someOption",
		Fields: []Field{
			{
				Name: "Name",
				Type: String,
			},
			{
				Name: "Flag",
				Type: Bool,
			},
		},
		Templates: []string{TemplateStruct, TemplateOptions},
		Variables: map[string]any{},
	}
	sb := &strings.Builder{}
	err := def.Execute(sb)
	if err != nil {
		t.Error(err)
		return
	}
	ct, err := format.Source([]byte(sb.String()))
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(string(ct))
}
