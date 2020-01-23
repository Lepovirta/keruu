package feeds2html

import (
	"html/template"
	"io"
)

const (
	templateStr = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta content="width=device-width,initial-scale=1" name="viewport">
  <title>{{ .Config.Title }}</title>
  <meta name="description" content="{{ .Config.Description }}">
  <meta property="og:title" content="{{ .Config.Title }}">
  <meta property="og:site_name" content="{{ .Config.Title }}">
  <meta property="og:type" content="website">
  <meta property="og:description" content="{{ .Config.Description }}">
  <style>{{ .Config.CSS }}</style>
</head>
<body>
<div class="wrapper">
  <header>
    <h1 class="title">{{ .Config.Title }}</h1>
    <p class="description">{{ .Config.Description }}</p>
  </header>
  <ul class="post-list">
  {{- range $count, $post := .Posts }}
  {{- with $post }}
    <li>
      <a class="post-title" href="{{ .Link }}">{{ .Title }}</a>
      <a class="post-feed-name" href="{{ .FeedLink }}">{{ .FeedName }}</a>
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

  defaultCSS = `body {
  background-color: #f1e8fc;
  padding: 0;
  margin: 0;
  font-family: sans-serif;
}

.wrapper {
  width: 700px;
  margin: 10px auto;
}

header {
  background-color: #44118d;
  color: #eee;
  padding: 5px 20px;
  margin: 0;
}

.title {
  margin: 2px 0;
}

.description {
  margin: 2px 0;
}

.post-list {
  background-color: white;
  margin: 0;
  padding: 20px;
}

.post-list li {
  list-style: none;
  padding: 5px 0;
}

.post-title {
  display: block;
  font-size: 1.2em;
}

.post-feed-name {
  font-size: 0.9em;
  color: #f44;
}

.post-time::before {
  content: '@ ';
}

.post-time {
  font-size: 0.9em;
  color: #484;
}

footer {
  border-top: 2px solid #eee;
  padding: 10px;
  font-size: 0.9em;
  text-align: right;
  color: #444;
}

@media screen and (max-width:700px) {
  .wrapper {
    width: auto;
    margin: 0;
  }
}
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

// renderHTML writes the given aggregation to the given writer in HTML format
func renderHTML(w io.Writer, a *aggregation) error {
	return htmlTemplate.Execute(w, a)
}
