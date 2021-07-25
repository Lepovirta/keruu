package feeds2html

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/mmcdole/gofeed"
)

type state struct {
	config      *Config
	feedParser  *gofeed.Parser
	aggregation *aggregation
	errorsFound int32
	mutex       sync.Mutex
}

func newState(config *Config) *state {
	feedParser := gofeed.NewParser()

	httpTranport := http.DefaultTransport.(*http.Transport).Clone()
	httpTranport.MaxIdleConns = 100
	httpTranport.MaxIdleConnsPerHost = 100
	httpTranport.MaxConnsPerHost = 100

	feedParser.Client = &http.Client{
		Transport: httpTranport,
		Timeout:   config.Fetch.HTTPTimeout,
	}

	aggregation := &aggregation{
		Config: &config.Aggregation,
		Posts:  make([]*feedPost, 0, 5000),
	}

	return &state{
		config:      config,
		feedParser:  feedParser,
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
	var wgFeeds sync.WaitGroup

	// Fetch feeds
	for _, feed := range s.config.Feeds {
		wgFeeds.Add(1)
		feed := feed
		go s.fetchFeed(&feed, &wgFeeds)
	}

	// Wait for everything to finish
	wgFeeds.Wait()

	// Post-process and write output
	s.postProcess()
	if err := s.writeHTML(out); err != nil {
		return err
	}

	// Check for errors
	if s.errorsFound > 0 {
		return fmt.Errorf("%d feed parsing errors found", s.errorsFound)
	}
	return nil
}

func (s *state) fetchFeed(feed *Feed, wg *sync.WaitGroup) {
	defer wg.Done()

	parsedFeed, err := s.feedParser.ParseURL(feed.URL.String())
	if err != nil {
		log.Printf("error processing feed '%s': %s", feed.URL, err)
		if s.config.Fetch.PropagateErrors {
			atomic.AddInt32(&s.errorsFound, 1)
		}
		return
	}

	posts := make([]*feedPost, 0, len(parsedFeed.Items))
	for _, item := range parsedFeed.Items {
		if isIncluded(feed, item.Title) {
			post := goFeedItemToPost(s.config, feed, parsedFeed, item)
			posts = append(posts, post)
		}
	}

	defer s.mutex.Unlock()
	s.mutex.Lock()
	s.aggregation.push(posts...)
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
