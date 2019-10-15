package feeds2html

import (
	"bufio"
	"github.com/Lepovirta/keruu/feeds2html/template"
	"github.com/mmcdole/gofeed"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Config struct {
	FeedURLs    []string
	HTTPTimeout time.Duration
	MaxItems    int
}

func DefaultConfig() *Config {
	return &Config{
		HTTPTimeout: time.Second * 10,
		MaxItems:    1000,
	}
}

type state struct {
	config     *Config
	feedParser *gofeed.Parser
	feedItemCh chan []*gofeed.Item
	feedItems  []*gofeed.Item
}

func newState(config *Config) *state {
	feedParser := gofeed.NewParser()
	feedParser.Client = &http.Client{
		Timeout: config.HTTPTimeout,
	}

	return &state{
		config:     config,
		feedParser: feedParser,
		feedItemCh: make(chan []*gofeed.Item, 100),
		feedItems:  make([]*gofeed.Item, 0, 5000),
	}
}

// Run fetches feeds specified in the config and outputs an HTML document for them
func Run(config *Config, out io.Writer) error {
	return newState(config).run(out)
}

func (s *state) run(out io.Writer) error {
	var wgJoiner sync.WaitGroup
	var wgFeeds sync.WaitGroup

	// Start feed item joiner
	wgJoiner.Add(1)
	go s.joinFeedItems(&wgJoiner)

	// Fetch feeds
	for _, url := range s.config.FeedURLs {
		wgFeeds.Add(1)
		go s.fetchFeed(url, &wgFeeds)
	}

	wgFeeds.Wait()
	close(s.feedItemCh)
	wgJoiner.Wait()
	s.postProcess()
	return s.writeHTML(out)
}

func (s *state) joinFeedItems(wg *sync.WaitGroup) {
	defer wg.Done()
	for items := range s.feedItemCh {
		s.feedItems = append(s.feedItems, items...)
	}
}

func (s *state) fetchFeed(url string, wg *sync.WaitGroup) {
	defer wg.Done()

	url = strings.TrimSpace(url)
	if url == "" {
		return
	}

	feed, err := s.feedParser.ParseURL(url)
	if err != nil {
		log.Printf("error processing feed '%s': %s", url, err)
	}
	s.feedItemCh <- feed.Items

	return
}

func (s *state) postProcess() {
	sortItemsByTime(s.feedItems)
	s.feedItems = s.feedItems[0:s.config.MaxItems]
}

func (s *state) writeHTML(w io.Writer) (err error) {
	bufw := bufio.NewWriter(w)

	err = template.Render(bufw, &template.Data{
		Items: s.feedItems,
	})
	if err != nil {
		return
	}

	err = bufw.Flush()
	return
}
