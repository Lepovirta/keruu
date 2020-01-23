package feeds2html

import (
	"sort"
	"time"

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
}

func goFeedItemToPost(
	feed *Feed,
	parsedFeed *gofeed.Feed,
	item *gofeed.Item,
) (post *feedPost, err error) {
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
