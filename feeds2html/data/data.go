package data

import (
	"html/template"
	"net/url"
	"sort"
	"time"

	"github.com/mmcdole/gofeed"
)

// Aggregation of posts from multiple feed sources
type Aggregation struct {
	Details *AggregationDetails
	Time    time.Time
	Posts   []*Post
}

// AggregationDetails contains information about the aggregation itself
type AggregationDetails struct {
	Title       string
	Description string
	MaxItems    int
	CSSString   string
}

// Post from a feed
type Post struct {
	Title string
	Link  *url.URL
	Time  *time.Time
}

// DefaultAggregationDetails provides default values for aggregation details
func DefaultAggregationDetails() *AggregationDetails {
	return &AggregationDetails{
		Title:       "Keruu",
		Description: "Aggregation of posts",
		MaxItems:    1000,
		CSSString:   defaultCSS,
	}
}

// GoFeedItemToPost converts an item from gofeed to a Post
func GoFeedItemToPost(item *gofeed.Item) (post *Post, err error) {
	link, err := url.Parse(item.Link)
	if err != nil {
		return
	}

	post = &Post{
		Title: item.Title,
		Link:  link,
		Time:  timeFromGoFeedItem(item),
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

// Push adds a new post to the aggregation
func (a *Aggregation) Push(post *Post) {
	a.Posts = append(a.Posts, post)
}

// Finalize performs the final post-processing for the aggregation data
func (a *Aggregation) Finalize() {
	sortPostsByTime(a.Posts)
	a.Posts = a.Posts[0:a.Details.MaxItems]
	a.Time = time.Now()
}

// FormattedTime provides the aggregation time a standard format
func (a *Aggregation) FormattedTime() string {
	return a.Time.Format("2006-01-02 15:04:05 -0700 MST")
}

// CSS provides the CSS data in template compatible format
func (ad *AggregationDetails) CSS() template.CSS {
	return template.CSS(ad.CSSString)
}

// Hostname provides the link hostname
func (p *Post) Hostname() string {
	return p.Link.Hostname()
}

// FormattedTime provides the post time in a standard string format
func (p *Post) FormattedTime() string {
	if p.Time == nil {
		return ""
	}
	return p.Time.Format("2006-01-02")
}

// After checks if this post is chronologically after the next one
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

func sortPostsByTime(posts []*Post) {
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].After(posts[j])
	})
}
