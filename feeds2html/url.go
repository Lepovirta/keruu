package feeds2html

import "net/url"

// URL wraps the standard library URL to provide extra functionality
type URL struct {
	*url.URL
}

// UnmarshalYAML parses an URL from a YAML formatted string
func (u *URL) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}
	url, err := url.Parse(s)
	if err != nil {
		return err
	}
	u.URL = url
	return nil
}
