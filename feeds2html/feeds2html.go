package feeds2html

import (
	"io"
	"bufio"
	"log"
	"github.com/mmcdole/gofeed"
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

func FromStream(feeds io.Reader, html io.Writer) error {
	bufHTML := bufio.NewWriter(html)
	if _, err := bufHTML.WriteString(htmlHeaderStr); err != nil {
		return err
	}

	feedScanner := bufio.NewScanner(feeds)
	for feedScanner.Scan() {
		feedURLStr := feedScanner.Text()
		if err := feedToHTML(feedURLStr, bufHTML); err != nil {
			log.Printf("error processing feed '%s': %s", feedURLStr, err)
		}
	}

	if err := feedScanner.Err(); err != nil {
		return err
	}

	if _, err := bufHTML.WriteString(htmlFooterStr); err != nil {
		return err
	}

	if err := bufHTML.Flush(); err != nil {
		return err
	}

	return nil
}

func feedToHTML(feedURLStr string, html *bufio.Writer) error {
	feedParser := gofeed.NewParser()
	feed, err := feedParser.ParseURL(feedURLStr)
	if err != nil {
		return err
	}

	htmlFeedTemplate.Execute(html, feed)

	return nil
}