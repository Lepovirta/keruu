package template

import (
	"html/template"
	"io"

	"github.com/Lepovirta/keruu/feeds2html/data"
)

const (
	templateStr = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta content="width=device-width,initial-scale=1" name="viewport">
  <title>{{ .Details.Title }}</title>
  <meta name="description" content="{{ .Details.Description }}">
  <meta property="og:title" content="{{ .Details.Title }}">
  <meta property="og:site_name" content="{{ .Details.Title }}">
  <meta property="og:type" content="website">
  <meta property="og:description" content="{{ .Details.Description }}">
  <style>{{ .Details.CSS }}</style>
</head>
<body>
<div class="wrapper">
  <header>
    <h1 class="title">{{ .Details.Title }}</h1>
    <p class="description">{{ .Details.Description }}</p>
  </header>
  <ul class="post-list">
  {{- range $count, $post := .Posts }}
  {{- with $post }}
    <li>
      <a class="post-title" href="{{ .Link }}">{{ .Title }}</a>
      <span class="post-hostname">{{ .Hostname }}</span>
      <span class="post-time">{{ .FormattedTime }}</span>
    </li>
  {{- end }}
  {{- end }}
  </ul>
  <footer>
    Generated using <a href="https://github.com/Lepovirta/keruu">keruu</a>
    at {{ .FormattedTime }}.
  </footer>
</div>
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

// Render writes the given aggregation to the given writer in HTML format
func Render(w io.Writer, aggregation *data.Aggregation) error {
	return htmlTemplate.Execute(w, aggregation)
}
