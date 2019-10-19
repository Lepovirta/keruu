package feeds2html

import (
	"fmt"
	"html/template"
	"io"
	"time"

	"gopkg.in/yaml.v2"
)

// Config contains the configuration for the entire feed fetching, aggregation, and rendering process
type Config struct {
	Feeds       []URL              `yaml:"feeds"`
	Fetch       *FetchConfig       `yaml:"fetch"`
	Aggregation *AggregationConfig `yaml:"aggregation"`
}

// FetchConfig contains the feed fetching related configurations
type FetchConfig struct {
	HTTPTimeout time.Duration `yaml:"httpTimeout"`
}

// AggregationConfig contains the feed aggregation related configurations
type AggregationConfig struct {
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	MaxItems    int    `yaml:"maxItems"`
	CSSString   string `yaml:"css"`
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

// ConfigFromYAML reads the configuration in YAML format
func ConfigFromYAML(r io.Reader) (c *Config, err error) {
	c = DefaultConfig()
	d := yaml.NewDecoder(r)
	err = d.Decode(c)
	return
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if len(c.Feeds) == 0 {
		return fmt.Errorf("no feeds provided")
	}
	if err := c.Fetch.Validate(); err != nil {
		return err
	}
	if err := c.Aggregation.Validate(); err != nil {
		return err
	}
	return nil
}

// Validate checks if the configuration is valid
func (c *FetchConfig) Validate() error {
	if c.HTTPTimeout <= 0 {
		return fmt.Errorf("HTTP timeout can't be zero")
	}
	return nil
}

// Validate checks if the configuration is valid
func (c *AggregationConfig) Validate() error {
	if c.MaxItems <= 0 {
		return fmt.Errorf("no point in limiting result size to 0")
	}
	return nil
}

// CSS provides the CSS data in HTML template compatible format
func (c *AggregationConfig) CSS() template.CSS {
	return template.CSS(c.CSSString)
}
