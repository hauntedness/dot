package dot

import (
	"io"
	"strings"
	"text/template"
)

// Struct which is going to be generated with [Templates] around it
type Struct struct {
	Name       string         // name of the struct
	TypeParams []TypeParam    // TypeParam
	Fields     Fields         // fields definition
	Directives Directives     // go directives of the struct
	Comments   Comments       // comments on top of struct
	Templates  Templates      // Templates with be rendered with a param map which hold other fields of Struct
	Variables  map[string]any // default nil, Variables store the additional values needed for the vendor template
}

// Execute execute template into writer
func (s *Struct) ExecuteTemplate(writer io.Writer, tmpl string) error {
	template, err := template.New("").Funcs(template.FuncMap{
		"Backquote": func(text string) string {
			return "`" + text + "`"
		},
	}).Parse(string(tmpl))
	if err != nil {
		return err
	}
	return template.Execute(writer, s)
}

// Execute execute template into string
func (s *Struct) ExecuteTemplateString(tmpl string) (string, error) {
	template, err := template.New("").Funcs(template.FuncMap{
		"Backquote": func(text string) string {
			return "`" + text + "`"
		},
	}).Parse(string(tmpl))
	if err != nil {
		return "", err
	}
	sb := &strings.Builder{}
	err = template.Execute(sb, s)
	if err != nil {
		return "", err
	}
	return sb.String(), nil
}

// Execute execute templates
func (s *Struct) Execute(writer io.Writer) error {
	for _, tmpl := range s.Templates {
		template, err := template.New("").Funcs(template.FuncMap{
			"Backquote": func(text string) string {
				return "`" + text + "`"
			},
		}).Parse(string(tmpl))
		if err != nil {
			return err
		}
		err = template.Execute(writer, s)
		if err != nil {
			return err
		}
	}
	return nil
}

// Field
type Field struct {
	Name      string         // field name
	Type      Type           // field type
	Tag       string         // tags
	Comments  Comments       // comments on top of the field
	Variables map[string]any // default nil, Variables store the additional values needed for the vendor template
}

// Fields
type Fields []Field

// Directives
type Directives []string

// Comments
type Comments []string

// Templates
type Templates []string

// Type
type Type struct {
	// represent the number of the pointer charactor "*"
	Stars int
	// Package
	Package Package
	Name    string // type name, can also be TypeParam name
}

// ID return the qualified type name
func (t Type) ID() string {
	result := ""
	if t.Package == (Package{}) {
		result = t.Name
	} else {
		result = t.Package.Name + t.Name
	}
	for i := 0; i < t.Stars; i++ {
		result = "*" + result
	}
	return result
}

// TypeParams
type TypeParams []TypeParam

// A TypeParam represents a type parameter type.
type TypeParam struct {
	TypeName   string // corresponding type name
	Constraint string
}

// Package
type Package struct {
	Path string
	Name string
}
