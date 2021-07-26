package aggregation

import (
	"io"
	"sort"
	"time"

	"gitlab.com/lepovirta/keruu/internal/feed"
)

type PostGroup struct {
	Name  string
	Posts []*feed.Post
}

type Aggregation struct {
	Config     *Config
	Time       time.Time
	PostGroups []PostGroup
}

func New(config *Config, posts []*feed.Post) *Aggregation {
	sortPostsByTime(posts)
	posts = posts[0:config.MaxPosts]

	return &Aggregation{
		Config:     config,
		Time:       time.Now(),
		PostGroups: groupPosts(posts, config.groupFunc()),
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
