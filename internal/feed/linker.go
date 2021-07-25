package feed

import (
	"net/url"
	"strings"

	"github.com/mmcdole/gofeed"
)

// ExtLink is a link to an external source
type ExtLink struct {
	Name string
	Link string
}

// Linker contains link patterns to other sites
type Linker struct {
	Name       string `yaml:"name"`
	URLPattern string `yaml:"url"`
}

func (l *Linker) goFeedItemToExtLink(item *gofeed.Item) ExtLink {
	link := strings.ReplaceAll(l.URLPattern, "$TITLE", url.QueryEscape(item.Title))
	link = strings.ReplaceAll(link, "$URL", url.QueryEscape(item.Link))
	return ExtLink{
		Name: l.Name,
		Link: link,
	}
}

func goFeedItemToExtLinks(linkers []Linker, item *gofeed.Item) []ExtLink {
	extLinks := make([]ExtLink, 0, len(linkers))
	for _, linker := range linkers {
		extLinks = append(extLinks, linker.goFeedItemToExtLink(item))
	}
	return extLinks
}
