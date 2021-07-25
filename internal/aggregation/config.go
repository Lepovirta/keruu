package aggregation

import (
	"fmt"
	"html/template"
)

// Config contains the feed aggregation related configurations
type Config struct {
	Title       string `yaml:"title,omitempty"`
	Description string `yaml:"description,omitempty"`
	MaxPosts    int    `yaml:"maxPosts,omitempty"`
	CSSString   string `yaml:"css,omitempty"`
}

func DefaultConfig() Config {
	return Config{
		Title:       "Keruu",
		Description: "Aggregation of posts",
		MaxPosts:    1000,
		CSSString:   defaultCSS,
	}
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.MaxPosts <= 0 {
		return fmt.Errorf("no point in limiting result size to 0")
	}
	return nil
}

// CSS provides the CSS data in HTML template compatible format
func (c *Config) CSS() template.CSS {
	return template.CSS(c.CSSString)
}
