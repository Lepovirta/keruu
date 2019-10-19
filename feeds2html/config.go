package feeds2html

import (
	"html/template"
	"time"
)

// Config contains the configuration for the entire feed fetching, aggregation, and rendering process
type Config struct {
	FeedURLs    []string
	Fetch       *FetchConfig
	Aggregation *AggregationConfig
}

// FetchConfig contains the feed fetching related configurations
type FetchConfig struct {
	HTTPTimeout time.Duration
}

// AggregationConfig contains the feed aggregation related configurations
type AggregationConfig struct {
	Title       string
	Description string
	MaxItems    int
	CSSString   string
}

// DefaultConfig generates a configuration with default values
func DefaultConfig() *Config {
	return &Config{
		Fetch: &FetchConfig{
			HTTPTimeout: time.Second * 10,
		},
		Aggregation: &AggregationConfig{
			Title:       "Keruu",
			Description: "Aggregation of posts",
			MaxItems:    1000,
			CSSString:   defaultCSS,
		},
	}
}

// CSS provides the CSS data in HTML template compatible format
func (ac *AggregationConfig) CSS() template.CSS {
	return template.CSS(ac.CSSString)
}
