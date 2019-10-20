package feeds2html

import (
	"bufio"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

// ConfigTemplateYAML is a YAML template for the config
const ConfigTemplateYAML = `feeds:
  - <URL>
  - <URL>
  - ...
fetch:
  httpTimeout: <DURATION>
aggregation:
  title: <STRING>
  description: <STRING>
  maxItems: <NUMBER>
  css: <STRING>
`

// Config contains the configuration for the entire feed fetching, aggregation, and rendering process
type Config struct {
	Feeds       []URL             `yaml:"feeds"`
	Fetch       FetchConfig       `yaml:"fetch,omitempty"`
	Aggregation AggregationConfig `yaml:"aggregation,omitempty"`
}

// FetchConfig contains the feed fetching related configurations
type FetchConfig struct {
	HTTPTimeout time.Duration `yaml:"httpTimeout,omitempty"`
}

// AggregationConfig contains the feed aggregation related configurations
type AggregationConfig struct {
	Title       string `yaml:"title,omitempty"`
	Description string `yaml:"description,omitempty"`
	MaxItems    int    `yaml:"maxItems,omitempty"`
	CSSString   string `yaml:"css,omitempty"`
}

// Init formats the config with default values
func (c *Config) Init() {
	c.Feeds = nil
	c.Fetch = FetchConfig{
		HTTPTimeout: time.Second * 10,
	}
	c.Aggregation = AggregationConfig{
		Title:       "Keruu",
		Description: "Aggregation of posts",
		MaxItems:    1000,
		CSSString:   defaultCSS,
	}
}

// FromYAMLFile reads the configuration from a YAML formatted file
func (c *Config) FromYAMLFile(filename string) (err error) {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Printf("failed to close input file: %s", err)
		}
	}()
	return c.FromYAML(bufio.NewReader(file))
}

// FromSTDIN reads the configuration from a STDIN in YAML format
func (c *Config) FromSTDIN() (err error) {
	return c.FromYAML(os.Stdin)
}

// FromYAML reads the configuration in YAML format
func (c *Config) FromYAML(r io.Reader) error {
	return yaml.NewDecoder(r).Decode(c)
}

// ToYAML converts the configuration to YAML format
func (c *Config) ToYAML(w io.Writer) error {
	return yaml.NewEncoder(w).Encode(c)
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
