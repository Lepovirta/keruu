package feeds2html

import (
	"bufio"
	"github.com/Lepovirta/keruu/feeds2html/template"
	"github.com/mmcdole/gofeed"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type Config struct {
	HTTPTimeout time.Duration
	MaxItems    int
}

func DefaultConfig() *Config {
	return &Config{
		HTTPTimeout: time.Second * 10,
		MaxItems:    1000,
	}
}

type Feeds2HTML struct {
	config     *Config
	feedParser *gofeed.Parser
	feedItems  []*gofeed.Item
}

func New(config *Config) *Feeds2HTML {
	feedParser := gofeed.NewParser()
	feedParser.Client = &http.Client{
		Timeout: config.HTTPTimeout,
	}

	return &Feeds2HTML{
		config:     config,
		feedParser: feedParser,
		feedItems:  make([]*gofeed.Item, 0, 5000),
	}
}

func (f *Feeds2HTML) FromStream(feeds io.Reader, out io.Writer) error {
	// Fetch feeds
	feedScanner := bufio.NewScanner(feeds)
	for feedScanner.Scan() {
		feedURLStr := feedScanner.Text()
		feedURLStr = strings.TrimSpace(feedURLStr)
		if feedURLStr == "" {
			break
		}
		if err := f.fetchItems(feedURLStr); err != nil {
			log.Printf("error processing feed '%s': %s", feedURLStr, err)
		}
	}

	if err := feedScanner.Err(); err != nil {
		log.Printf("error reading feed links: %s", err)
	}

	// Sort and slice
	sortItemsByTime(f.feedItems)
	f.feedItems = f.feedItems[0:f.config.MaxItems]

	// Write output
	return f.writeHTML(out)
}

func (f *Feeds2HTML) fetchItems(feedURLStr string) error {
	feed, err := f.feedParser.ParseURL(feedURLStr)
	if err != nil {
		return err
	}
	f.feedItems = append(f.feedItems, feed.Items...)
	return nil
}

func (f *Feeds2HTML) writeHTML(w io.Writer) (err error) {
	bufw := bufio.NewWriter(w)

	err = template.Render(bufw, &template.Data{
		Items: f.feedItems,
	})
	if err != nil {
		return
	}

	err = bufw.Flush()
	return
}
