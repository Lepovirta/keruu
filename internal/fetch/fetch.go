package fetch

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/mmcdole/gofeed"
	"gitlab.com/lepovirta/keruu/internal/feed"
)

type Config struct {
	HTTPTimeout     time.Duration `yaml:"httpTimeout,omitempty"`
	PropagateErrors bool          `yaml:"propagateErrors"`
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.HTTPTimeout <= 0 {
		return fmt.Errorf("HTTP timeout can't be zero")
	}
	return nil
}

func DefaultConfig() Config {
	return Config{
		HTTPTimeout: time.Second * 10,
	}
}

type state struct {
	config         *Config
	feeds          []feed.Config
	linkers        []feed.Linker
	feedParser     *gofeed.Parser
	collectedPosts []*feed.Post
	errorsFound    int32
	mutex          sync.Mutex
}

func newState(config *Config, feeds []feed.Config, linkers []feed.Linker) *state {
	feedParser := gofeed.NewParser()

	httpTranport := http.DefaultTransport.(*http.Transport).Clone()
	httpTranport.MaxIdleConns = 100
	httpTranport.MaxIdleConnsPerHost = 100
	httpTranport.MaxConnsPerHost = 100

	feedParser.Client = &http.Client{
		Transport: httpTranport,
		Timeout:   config.HTTPTimeout,
	}

	return &state{
		config:         config,
		feeds:          feeds,
		linkers:        linkers,
		feedParser:     feedParser,
		collectedPosts: make([]*feed.Post, 0, 5000),
	}
}

func (s *state) run() error {
	var wg sync.WaitGroup

	// Fetch feeds
	for _, f := range s.feeds {
		wg.Add(1)
		f := f
		go s.fetchFeed(&f, &wg)
	}

	// Wait for everything to finish
	wg.Wait()

	// Check for errors
	if s.errorsFound > 0 {
		return fmt.Errorf("%d feed parsing errors found", s.errorsFound)
	}
	return nil
}

func (s *state) fetchFeed(f *feed.Config, wg *sync.WaitGroup) {
	defer wg.Done()

	parsedFeed, err := s.feedParser.ParseURL(f.URL.String())
	if err != nil {
		log.Printf("error processing feed '%s': %s", f.URL, err)
		if s.config.PropagateErrors {
			atomic.AddInt32(&s.errorsFound, 1)
		}
		return
	}

	posts := make([]*feed.Post, 0, len(parsedFeed.Items))
	for _, item := range parsedFeed.Items {
		if f.IsIncluded(item.Title) {
			post := f.PostFromGoFeedItem(s.linkers, parsedFeed, item)
			posts = append(posts, post)
		}
	}

	defer s.mutex.Unlock()
	s.mutex.Lock()
	s.collectedPosts = append(s.collectedPosts, posts...)
}

func Run(config *Config, feeds []feed.Config, linkers []feed.Linker) ([]*feed.Post, error) {
	state := newState(config, feeds, linkers)
	err := state.run()
	return state.collectedPosts, err
}
