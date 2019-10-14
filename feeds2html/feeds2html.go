package feeds2html

import (
	"bufio"
	"github.com/mmcdole/gofeed"
	"io"
	"log"
	"net/http"
	"time"
)

type Config struct {
	HTTPTimeout time.Duration
}

func DefaultConfig() *Config {
	return &Config{
		HTTPTimeout: time.Second * 10,
	}
}

type Feeds2HTML struct {
	feedParser *gofeed.Parser
}

func New(config *Config) *Feeds2HTML {
	feedParser := gofeed.NewParser()
	feedParser.Client = &http.Client{
		Timeout: config.HTTPTimeout,
	}

	return &Feeds2HTML{
		feedParser: feedParser,
	}
}

func (f *Feeds2HTML) FromStream(feeds io.Reader, html io.Writer) error {
	bufHTML := bufio.NewWriter(html)
	if _, err := bufHTML.WriteString(htmlHeaderStr); err != nil {
		return err
	}

	feedScanner := bufio.NewScanner(feeds)
	for feedScanner.Scan() {
		feedURLStr := feedScanner.Text()
		if err := f.feedToHTML(feedURLStr, bufHTML); err != nil {
			log.Printf("error processing feed '%s': %s", feedURLStr, err)
		}
	}

	if err := feedScanner.Err(); err != nil {
		return err
	}

	if _, err := bufHTML.WriteString(htmlFooterStr); err != nil {
		return err
	}

	if err := bufHTML.Flush(); err != nil {
		return err
	}

	return nil
}

func (f *Feeds2HTML) feedToHTML(feedURLStr string, html *bufio.Writer) error {
	feed, err := f.feedParser.ParseURL(feedURLStr)
	if err != nil {
		return err
	}

	htmlFeedTemplate.Execute(html, feed)

	return nil
}
