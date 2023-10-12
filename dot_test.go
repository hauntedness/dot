package dot

import (
	"fmt"
	"go/format"
	"strings"
	"testing"
	"text/template"
)

func TestStruct_ExecuteTemplateString(t *testing.T) {
	st := Struct{
		Name: "TestStruct",
		TypeParams: []TypeParam{
			{
				TypeName:   "T",
				Constraint: "any",
			},
			{
				TypeName:   "R",
				Constraint: "string",
			},
		},
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

func TestStruct_Quote(t *testing.T) {
	tmpl, err := template.New("").Funcs(template.FuncMap{
		"Backquote": func(text string) string {
			return "`" + text + "`"
		},
	}).Parse(`{{ range .TypeParams }} {{index .}} {{else}} T0  {{end}}`)
	if err != nil {
		t.Error(err)
		return
	}
	st := Struct{
		TypeParams: []TypeParam{
			{
				TypeName:   "T",
				Constraint: "any",
			},
		},
		Variables: map[string]any{"1": 1, "2": 2, "3": 3},
	}
	sb := &strings.Builder{}
	_ = tmpl.Execute(sb, st)
	fmt.Println(sb.String())
}
