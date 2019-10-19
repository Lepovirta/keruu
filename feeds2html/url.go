package feeds2html

import "net/url"

type URL struct {
	*url.URL
}

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
