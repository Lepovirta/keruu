package feed

import (
	"time"

	"github.com/mmcdole/gofeed"
	"gitlab.com/lepovirta/keruu/internal/regexp"
	"gitlab.com/lepovirta/keruu/internal/url"
)

// Config contains the details of a single feed
type Config struct {
	Name    string      `yaml:"name"`
	URL     url.URL     `yaml:"url"`
	Exclude []regexp.RE `yaml:"exclude,omitempty"`
	Include []regexp.RE `yaml:"include,omitempty"`
}

func (f *Config) IsIncluded(s string) bool {
	// Always exclude first
	for _, filter := range f.Exclude {
		if filter.MatchString(s) {
			return false
		}
	}

	for _, filter := range f.Include {
		if filter.MatchString(s) {
			return true
		}
	}

	// No match on filters?
	// Only include it, if there was no include filters set
	return len(f.Include) == 0
}

func (f *Config) PostFromGoFeedItem(
	linkers []Linker,
	parsedFeed *gofeed.Feed,
	item *gofeed.Item,
) *Post {
	feedName := f.Name
	if feedName == "" {
		feedName = parsedFeed.Title
	}

	return &Post{
		FeedName: feedName,
		FeedLink: parsedFeed.Link,
		Title:    item.Title,
		Link:     item.Link,
		Time:     timeFromGoFeedItem(item),
		ExtLinks: goFeedItemToExtLinks(linkers, item),
	}
}

type Post struct {
	FeedName string
	FeedLink string
	Title    string
	Link     string
	Time     *time.Time
	ExtLinks []ExtLink
}

func timeFromGoFeedItem(item *gofeed.Item) *time.Time {
	if item.PublishedParsed != nil {
		return item.PublishedParsed
	} else if item.UpdatedParsed != nil {
		return item.UpdatedParsed
	}
	return nil
}

func (p *Post) FormattedTime() string {
	if p.Time == nil {
		return ""
	}
	return p.Time.Format("2006-01-02")
}

func (p *Post) After(other *Post) bool {
	// No timestamp means it's considered older
	if p.Time == nil {
		return false
	} else if other.Time == nil {
		return true
	}

	// Newest first
	return p.Time.After(*other.Time)
}