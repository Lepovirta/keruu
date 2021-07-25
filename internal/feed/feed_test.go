package feed

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/lepovirta/keruu/internal/regexp"
)

func TestIsIncludedExcludeOnly(t *testing.T) {
	feed := &Config{
		Exclude: []regexp.RE{
			regexp.MustCompile("^Sponsored Post:"),
			regexp.MustCompile("\\[ad\\]$"),
		},
	}

	assert.True(t, feed.IsIncluded("I'm just a normal post title"))
	assert.False(t, feed.IsIncluded("Sponsored Post: Get cryptobux for free!"))
	assert.False(t, feed.IsIncluded("5 reasons why CorpCloud is the best [ad]"))
}

func TestIsIncludedIncludeOnly(t *testing.T) {
	feed := &Config{
		Include: []regexp.RE{
			regexp.MustCompile("^Podcast:"),
		},
	}

	assert.True(t, feed.IsIncluded("Podcast: Kubernetes in Space"))
	assert.False(t, feed.IsIncluded("Monthly newsletter"))
}

func TestIsIncludedMixed(t *testing.T) {
	feed := &Config{
		Exclude: []regexp.RE{
			regexp.MustCompile("\\[ad\\]$"),
		},
		Include: []regexp.RE{
			regexp.MustCompile("^Podcast:"),
		},
	}

	assert.True(t, feed.IsIncluded("Podcast: Kubernetes in Space"))
	assert.False(t, feed.IsIncluded("Podcast: 5 reasons why CorpCloud is the best [ad]"))
	assert.False(t, feed.IsIncluded("Monthly newsletter"))
}
