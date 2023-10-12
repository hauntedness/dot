package dot

import (
	"fmt"
	"go/format"
	"testing"
)

func TestStruct_ExecuteTemplateString(t *testing.T) {
	st := Struct{
		Name:       "TestStruct",
		TypeParams: "T any, R string",
		Fields: []Field{
			{
				Name:     "Name",
				Type:     Type{Name: "string"},
				Tag:      `json:"name"`,
				Comments: []string{"the name"},
			},
			{
				Name:     "Age",
				Type:     Type{Name: "int"},
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
