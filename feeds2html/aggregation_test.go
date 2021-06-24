package feeds2html

import (
	"regexp"
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

func reFromString(s string) RE {
	re := regexp.MustCompile(s)
	return NewRE(re)
}

func TestIsIncludedExcludeOnly(t *testing.T) {
	feed := &Feed{
		Exclude: []RE{
			reFromString("^Sponsored Post:"),
			reFromString("\\[ad\\]$"),
		},
	}

	assert.True(t, isIncluded(feed, "I'm just a normal post title"))
	assert.False(t, isIncluded(feed, "Sponsored Post: Get cryptobux for free!"))
	assert.False(t, isIncluded(feed, "5 reasons why CorpCloud is the best [ad]"))
}

func TestIsIncludedIncludeOnly(t *testing.T) {
	feed := &Feed{
		Include: []RE{
			reFromString("^Podcast:"),
		},
	}

	assert.True(t, isIncluded(feed, "Podcast: Kubernetes in Space"))
	assert.False(t, isIncluded(feed, "Monthly newsletter"))
}

func TestIsIncludedMixed(t *testing.T) {
	feed := &Feed{
		Exclude: []RE{
			reFromString("\\[ad\\]$"),
		},
		Include: []RE{
			reFromString("^Podcast:"),
		},
	}

	assert.True(t, isIncluded(feed, "Podcast: Kubernetes in Space"))
	assert.False(t, isIncluded(feed, "Podcast: 5 reasons why CorpCloud is the best [ad]"))
	assert.False(t, isIncluded(feed, "Monthly newsletter"))
}
