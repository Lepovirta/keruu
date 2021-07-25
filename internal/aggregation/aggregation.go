package aggregation

import (
	"io"
	"sort"
	"time"

	"gitlab.com/lepovirta/keruu/internal/feed"
)

type Aggregation struct {
	Config *Config
	Time   time.Time
	Posts  []*feed.Post
}

func New(config *Config, posts []*feed.Post) *Aggregation {
	sortPostsByTime(posts)
	return &Aggregation{
		Config: config,
		Time: time.Now(),
		Posts: posts[0:config.MaxPosts],
	}
}

func sortPostsByTime(posts []*feed.Post) {
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].After(posts[j])
	})
}

func (a *Aggregation) FormattedTime() string {
	return a.Time.Format("2006-01-02 15:04:05 -0700 MST")
}

func (a *Aggregation) ToHTML(w io.Writer) error {
	return renderHTML(w, a)
}
