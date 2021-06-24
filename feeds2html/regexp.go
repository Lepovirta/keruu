package feeds2html

import (
	"regexp"
)

// RE wraps the standard library Regexp to provide extra functionality.
// Specifically, it provides YAML unmarshalling.
type RE struct {
	*regexp.Regexp
}

// NewRE converts Go Regexp to our custom RE type
func NewRE(re *regexp.Regexp) RE {
	return RE{re}
}

// UnmarshalYAML parses an regular expression from a YAML formatted string
func (r *RE) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}
	pattern, err := regexp.Compile(s)
	if err != nil {
		return err
	}
	r.Regexp = pattern
	return nil
}
