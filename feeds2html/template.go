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
  <link rel="icon" href="data:,">
  <style>{{ .Config.CSS }}</style>
</head>
<body>
<main>
  <header>
    <h1 class="title">{{ .Config.Title }}</h1>
    <p class="description">{{ .Config.Description }}</p>
  </header>
  <ul class="post-list">
  {{- range $_, $post := .Posts }}
  {{- with $post }}
    <li>
      <div class="post-head">
        <a class="post-title" href="{{ .Link }}">{{ .Title }}</a>
      </div>
      <div class="post-meta">
        <a class="post-feed" href="{{ .FeedLink }}">{{ .FeedName }}</a>
        <time class="post-time" datetime="{{ .FormattedTime }}">{{ .FormattedTime }}</time>
        {{- if .ExtLinks }}
        <ul class="post-ext-links">
        {{- range $_, $extLink := .ExtLinks }}
          <li><a href="{{ $extLink.Link }}">{{ $extLink.Name }}</a></li>
        {{- end }}
        </ul>
        {{- end }}
      </div>
    </li>
  {{- end }}
  {{- end }}
  </ul>
  <footer>
    Generated using <a href="https://gitlab.com/lepovirta/keruu">keruu</a>
    at {{ .FormattedTime }}.
  </footer>
</main>
</body>
</html>
`

	defaultCSS = `body {
  background-color: #f7f0ff;
  padding: 0;
  margin: 0;
  font-family: sans-serif;
}

main {
  max-width: 760px;
  margin: 20px auto;
}

header {
  background-color: #44118d;
  color: #eee;
  padding: 5px 20px;
  margin: 0;
}

a {
  color: #00e;
}

a:visited {
  color: #551abb;
}

footer {
  padding: 10px;
  font-size: 0.9em;
  text-align: right;
  color: #444;
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

.post-list a {
  text-decoration: none;
}

.post-list a:hover {
  text-decoration: underline;
}

.post-list li {
  list-style: none;
  padding: 5px 0;
}

.post-head {
  font-size: 1.2em;
}

.post-meta {
  font-size: 0.9em;
}

.post-count::after {
  content: '.';
}

.post-feed {
  color: #f44;
  text-decoration: none;
}

.post-feed-name:visited {
  color: #b22;
}

.post-time::before {
  content: '| ';
}

.post-time {
  font-size: 0.9em;
  color: #484;
}

.post-ext-links {
  display: inline;
  padding: 0;
}

.post-ext-links::before {
  content: '|';
  color: #18798e;
}

.post-ext-links li {
  display: inline;
}

.post-ext-links li::after {
  content: ',';
}

.post-ext-links li:nth-last-child(1)::after {
  content: '';
}

.post-ext-links a {
  color: #18798e;
}

.post-ext-links a:visited {
  color: #104754;
}

@media screen and (max-width:760px) {
  main {
    max-width: auto;
    margin: 0;
  }
}

@media (prefers-color-scheme: dark) {
  body {
    background-color: #113;
  }
  a {
    color: #aef;
  }
  a:visited {
    color: #8be;
  }
  footer {
    color: #baf;
  }
  .post-list {
    background-color: #214;
  }
  .post-feed-name {
    color: #f88;
  }
  .post-feed-name:visited {
    color: #b66;
  }
  .post-time {
    color: #8f8;
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
