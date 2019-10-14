package feeds2html

import (
	"html/template"
)

const (
	htmlHeaderStr = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>Keruu</title>
  <meta name="description" content="Keruu">
</head>
<body>
<ul>`
	htmlFooterStr = `</ul>
</body>
</html>
`
	htmlFeedTemplateStr = `
{{- range $count, $item := .Items }}
<li><a href="{{ .Link }}">{{ .Title }}</a></li>
{{- end }}
`
)

var (
	htmlFeedTemplate *template.Template
)

func init() {
	var err error

	htmlFeedTemplate, err = template.New("feed").Parse(htmlFeedTemplateStr)
	if err != nil {
		panic(err)
	}
}
