package dot

const TemplateStruct = `
{{- define "TemplateTypeParams" }}[
	{{- if .TypeParams -}} 
		{{- .TypeParams}}, 
	{{- end -}}]
{{- end}}
{{- range .Comments }}
	{{- "// " }} {{- . }}
{{- end }}
	//
{{- range .Directives }}
	// {{- . }}
{{- end }}
type {{.Name}} {{template "TemplateTypeParams" .}} struct{
	{{ range .Fields }}
		{{- range .Comments -}}
			// {{- . }}
		{{ end }}
		{{- if .Tag}}
		{{- Backquote .Tag | printf "%s %s %s" .Name .Type.FullName }}
		{{- else}}
		{{- printf "%s %s" .Name .Type.FullName }}
		{{- end}}
	{{ end }}
}
`

const TemplateOptions = `
{{$struct := .}}
{{ range .Fields }}
var With{{ Transform .Name "ToUpper1st" }} = func({{ Transform .Name "ToLower1st" }} {{ .Type.FullName }}) func(*{{$struct.Name}}) {
	return func(o *{{$struct.Name}}) {
		o.{{.Name}} = {{ Transform .Name "ToLower1st" }}
	}
}
{{ end }}
`
