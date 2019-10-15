package template

import (
	"github.com/mmcdole/gofeed"
	"html/template"
	"io"
)

const (
	templateStr = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>Keruu</title>
  <meta name="description" content="Keruu">
</head>
<body>
<ul>
{{- range $count, $item := .Items }}
{{- with $item }}
<li><a href="{{ .Link }}">{{ .Title }}</a></li>
{{- end }}
{{- end }}
</ul>
</body>
</html>
`
)

var (
	htmlTemplate *template.Template
)

func init() {
	var err error

	htmlTemplate, err = template.New("").Parse(templateStr)
	if err != nil {
		panic(err)
	}
}

type Data struct {
	Items []*gofeed.Item
}

func Render(w io.Writer, data *Data) error {
	return htmlTemplate.Execute(w, data)
}
