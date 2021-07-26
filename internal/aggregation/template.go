package aggregation

import (
	_ "embed"
	"html/template"
	"io"
)

var (
	htmlTemplate *template.Template

	//go:embed style.css
	defaultCSS string

	//go:embed template.html
	htmlTemplateStr string
)

func init() {
	var err error

	htmlTemplate, err = template.New("").Parse(htmlTemplateStr)
	if err != nil {
		panic(err)
	}
}

// renderHTML writes the given aggregation to the given writer in HTML format
func renderHTML(w io.Writer, a *Aggregation) error {
	return htmlTemplate.Execute(w, a)
}
