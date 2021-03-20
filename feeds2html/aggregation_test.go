package feeds2html

import (
	"testing"

	"github.com/mmcdole/gofeed"
	"github.com/stretchr/testify/assert"
)

func TestExtLinkPatternToURL(t *testing.T) {
	linker1 := Linker{
		Name: "DuckDuckGo",
		URLPattern: "https://duckduckgo.com/?q=$TITLE",
	}
	linker2 := Linker{
		Name: "reddit",
		URLPattern: "https://old.reddit.com/submit?url=$URL",
	}
	post := gofeed.Item{
		Title: "Hello World",
		Link: "http://example.org",
	}

	url1 := extLinkPatternToURL(&linker1, &post)
	url2 := extLinkPatternToURL(&linker2, &post)

	assert.Equal(t, "https://duckduckgo.com/?q=Hello+World", url1)
	assert.Equal(t, "https://old.reddit.com/submit?url=http%3A%2F%2Fexample.org", url2)
}
