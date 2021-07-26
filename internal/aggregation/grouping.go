package aggregation

import (
	"fmt"
	"strconv"

	"gitlab.com/lepovirta/keruu/internal/feed"
)

const (
	weeklyGrouping  string = "weekly"
	monthlyGrouping        = "monthly"
	noGrouping             = "none"
	defaultGrouping        = monthlyGrouping
)

type GroupFunc func(index int, post *feed.Post) string

func weeklyGroupF(index int, post *feed.Post) string {
	year, week := post.Time.ISOWeek()
	return fmt.Sprintf("%d/%d", week, year)
}

func monthlyGroupF(index int, post *feed.Post) string {
	return post.Time.Format("01/2006")
}

func everyNthGroupF(n int) GroupFunc {
	return func(index int, post *feed.Post) string {
		return fmt.Sprint(index/n + 1)
	}
}

func groupingStringToFunc(s string) GroupFunc {
	i, err := strconv.Atoi(s)
	if err == nil {
		return everyNthGroupF(int(i))
	}

	switch s {
	case weeklyGrouping:
		return weeklyGroupF
	case monthlyGrouping:
		return monthlyGroupF
	case noGrouping, "":
		return nil
	default:
		return nil
	}
}

func isValidGrouping(s string) bool {
	i, err := strconv.Atoi(s)
	if err == nil && i > 0 {
		return true
	}

	switch s {
	case weeklyGrouping, monthlyGrouping, noGrouping, "":
		return true
	default:
		return false
	}
}

func groupPosts(posts []*feed.Post, groupF GroupFunc) []PostGroup {
	// No grouping specified => all in one unamed group
	if groupF == nil {
		return []PostGroup{
			{
				Name:  "",
				Posts: posts,
			},
		}
	}

	groupIndex := 0
	groups := make([]PostGroup, 0, 60)
	groups = append(groups, PostGroup{})

	for index, post := range posts {
		name := groupF(index, post)
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
