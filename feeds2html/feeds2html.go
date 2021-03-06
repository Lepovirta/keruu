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

	httpTranport := http.DefaultTransport.(*http.Transport).Clone()
	httpTranport.MaxIdleConns = 100
	httpTranport.MaxIdleConnsPerHost = 100
	httpTranport.MaxConnsPerHost = 100

	feedParser.Client = &http.Client{
		Transport: httpTranport,
		Timeout: config.Fetch.HTTPTimeout,
	}

	aggregation := &aggregation{
		Config: &config.Aggregation,
		Posts:  make([]*feedPost, 0, 5000),
	}

	return &state{
		config:      config,
		feedParser:  feedParser,
		postCh:      make(chan *feedPost, 1000),
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
	for _, feed := range s.config.Feeds {
		wgFeeds.Add(1)
		feed := feed
		go s.fetchFeed(&feed, &wgFeeds)
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

func (s *state) fetchFeed(feed *Feed, wg *sync.WaitGroup) {
	defer wg.Done()

	parsedFeed, err := s.feedParser.ParseURL(feed.URL.String())
	if err != nil {
		log.Printf("error processing feed '%s': %s", feed.URL, err)
		return
	}
	for _, item := range parsedFeed.Items {
		if isIncluded(feed, item.Title) {
			post := goFeedItemToPost(s.config, feed, parsedFeed, item)
			s.postCh <- post
		}
	}
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
