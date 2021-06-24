package feeds2html

import (
	"sort"
	"strings"
	"time"
	"net/url"

	"github.com/mmcdole/gofeed"
)

type aggregation struct {
	Config *AggregationConfig
	Time   time.Time
	Posts  []*feedPost
}

type feedPost struct {
	FeedName string
	FeedLink string
	Title    string
	Link     string
	Time     *time.Time
	ExtLinks  []extLink
}

type extLink struct {
	Name string
	Link string
}

func isIncluded(feed *Feed, s string) bool {
	// Always exclude first
	for _, filter := range feed.Exclude {
		if filter.MatchString(s) {
			return false
		}
	}

	for _, filter := range feed.Include {
		if filter.MatchString(s) {
			return true
		}
	}

	// No match on filters?
	// Only include it, if there was no include filters set
	return len(feed.Include) == 0
}

func goFeedItemToPost(
	config *Config,
	feed *Feed,
	parsedFeed *gofeed.Feed,
	item *gofeed.Item,
) (post *feedPost) {
	feedName := feed.Name
	if feedName == "" {
		feedName = parsedFeed.Title
	}

	post = &feedPost{
		FeedName: feedName,
		FeedLink: parsedFeed.Link,
		Title:    item.Title,
		Link:     item.Link,
		Time:     timeFromGoFeedItem(item),
		ExtLinks: extLinksFromGoFeedItem(config, item),
	}

	return
}

func timeFromGoFeedItem(item *gofeed.Item) *time.Time {
	if item.PublishedParsed != nil {
		return item.PublishedParsed
	} else if item.UpdatedParsed != nil {
		return item.UpdatedParsed
	}
	return nil
}

func extLinksFromGoFeedItem(config *Config, item *gofeed.Item) []extLink {
	extLinks := make([]extLink, 0, len(config.Links))
	for _, linker := range config.Links {
		extLinks = append(extLinks, extLink{
			Name: linker.Name,
			Link: extLinkPatternToURL(&linker, item),
		})
	}
	return extLinks
}

func extLinkPatternToURL(linker *Linker, item *gofeed.Item) (u string) {
	u = strings.ReplaceAll(linker.URLPattern, "$TITLE", url.QueryEscape(item.Title))
	u = strings.ReplaceAll(u, "$URL", url.QueryEscape(item.Link))
	return
}

func (a *aggregation) push(post *feedPost) {
	a.Posts = append(a.Posts, post)
}

func (a *aggregation) finalize() {
	sortPostsByTime(a.Posts)
	a.Posts = a.Posts[0:a.Config.MaxPosts]
	a.Time = time.Now()
}

func (a *aggregation) FormattedTime() string {
	return a.Time.Format("2006-01-02 15:04:05 -0700 MST")
}

func (p *feedPost) FormattedTime() string {
	if p.Time == nil {
		return ""
	}
	return p.Time.Format("2006-01-02")
}

func (p *feedPost) after(other *feedPost) bool {
	// No timestamp means it's considered older
	if p.Time == nil {
		return false
	} else if other.Time == nil {
		return true
	}

	// Newest first
	return p.Time.After(*other.Time)
}

func sortPostsByTime(posts []*feedPost) {
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].after(posts[j])
	})
}
