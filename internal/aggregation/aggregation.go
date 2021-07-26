package aggregation

import (
	"fmt"
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
		PostGroups: monthlyGroups(posts),
	}
}

func groupPosts(posts []*feed.Post, groupName func(*feed.Post) string) []PostGroup {
	groupIndex := 0
	groups := make([]PostGroup, 0, 60)
	groups = append(groups, PostGroup{})

	for _, post := range posts {
		name := groupName(post)
		group := &groups[groupIndex]

		// If the current group name doesn't match the computed name,
		// create a new group.
		if group.Name != "" && group.Name != name {
			groupIndex += 1
			groups = append(groups, PostGroup{})
			group = &groups[groupIndex]
		}

		// Initialize the group if needed
		if group.Name == "" {
			group.Name = name
			group.Posts = make([]*feed.Post, 100)
		}

		// Finally, add the post to the current group
		group.Posts = append(group.Posts, post)
	}
	return groups
}

func monthlyGroups(posts []*feed.Post) []PostGroup {
	return groupPosts(posts, func(p *feed.Post) string {
		return p.Time.Format("01/2006")
	})
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
