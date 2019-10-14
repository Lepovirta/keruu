package feeds2html

import (
	"github.com/mmcdole/gofeed"
	"sort"
)

func sortItemsByTime(items []*gofeed.Item) {
	sort.Slice(items, func(i, j int) bool {
		// Check either published or updated time from the first candidate
		first := items[i].PublishedParsed
		if first == nil {
			first = items[i].UpdatedParsed
		}

		// Check either published or updated time from the second candidate
		second := items[j].PublishedParsed
		if second == nil {
			second = items[j].UpdatedParsed
		}

		// No timestamp means it's considered old
		if first == nil {
			return false
		} else if second == nil {
			return true
		}

		// Newest first
		return first.After(*second)
	})
}
