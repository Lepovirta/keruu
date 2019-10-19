package feeds2html

import (
	"bufio"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/mmcdole/gofeed"
)

type state struct {
	config      *Config
	feedParser  *gofeed.Parser
	postCh      chan *feedPost
	aggregation *aggregation
}

func newState(config *Config) *state {
	feedParser := gofeed.NewParser()
	feedParser.Client = &http.Client{
		Timeout: config.Fetch.HTTPTimeout,
	}

	aggregation := &aggregation{
		Config: config.Aggregation,
		Posts:  make([]*feedPost, 0, 5000),
	}

	return &state{
		config:      config,
		feedParser:  feedParser,
		postCh:      make(chan *feedPost, 100),
		aggregation: aggregation,
	}
}

// Run fetches feeds specified in the config and outputs an HTML document for them
func Run(config *Config, out io.Writer) error {
	if err := config.Validate(); err != nil {
		return err
	}
	return newState(config).run(out)
}

func (s *state) run(out io.Writer) error {
	var wgJoiner sync.WaitGroup
	var wgFeeds sync.WaitGroup

	// Start feed item joiner
	wgJoiner.Add(1)
	go s.joinFeedItems(&wgJoiner)

	// Fetch feeds
	for _, url := range s.config.Feeds {
		wgFeeds.Add(1)
		go s.fetchFeed(url, &wgFeeds)
	}

	// Wait for everything to finish
	wgFeeds.Wait()
	close(s.postCh)
	wgJoiner.Wait()

	// Post-process and write output
	s.postProcess()
	return s.writeHTML(out)
}

func (s *state) joinFeedItems(wg *sync.WaitGroup) {
	defer wg.Done()
	for post := range s.postCh {
		s.aggregation.push(post)
	}
}

func (s *state) fetchFeed(url URL, wg *sync.WaitGroup) {
	defer wg.Done()

	feed, err := s.feedParser.ParseURL(url.String())
	if err != nil {
		log.Printf("error processing feed '%s': %s", url, err)
		return
	}
	for _, item := range feed.Items {
		post, err := goFeedItemToPost(item)
		if err != nil {
			log.Printf("error processing post: %s", err)
			break
		}
		s.postCh <- post
	}

	return
}

func (s *state) postProcess() {
	s.aggregation.finalize()
}

func (s *state) writeHTML(w io.Writer) (err error) {
	bufw := bufio.NewWriter(w)

	err = renderHTML(bufw, s.aggregation)
	if err != nil {
		return
	}

	err = bufw.Flush()
	return
}
