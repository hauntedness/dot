package dot

var TemplateStruct = `
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
		{{- Backquote .Tag | printf "%s %s %s" .Name .Type.ID }}
	{{ end }}
}
`
