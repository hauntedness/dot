package dot

import (
	"strings"
	"text/template"
	"unicode"
)

var (
	Int     = NewType("int")
	Uint    = NewType("uint")
	Float32 = NewType("float32")
	Float64 = NewType("float64")
	Bool    = NewType("bool")
	String  = NewType("string")
)

var FuncMap = template.FuncMap{
	"Backquote": func(text string) string {
		return "`" + text + "`"
	},
	"Transform": func(text string, indicator string) string {
		if len(text) == 0 {
			return text
		}
		runes := []rune(text)
		if strings.EqualFold(indicator, "ToUpper1st") {
			runes[0] = unicode.ToUpper(runes[0])
		} else if strings.EqualFold(indicator, "ToLower1st") {
			runes[0] = unicode.ToLower(runes[0])
		}
		return string(runes)
	},
}
