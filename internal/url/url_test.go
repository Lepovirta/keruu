package url

import (
	"net/url"
	"testing"

	"gopkg.in/yaml.v3"
	"github.com/stretchr/testify/assert"
)

const urlYAML = `
- http://example.org/
- https://lepovirta.org/posts/index.html
- https://duckduckgo.com/?q=helloworld
`

var urlList = []string{
	"http://example.org/",
	"https://lepovirta.org/posts/index.html",
	"https://duckduckgo.com/?q=helloworld",
}

func TestURLUnmarshallingToYAML(t *testing.T) {
	// setup expected URLs
	expectedURLs := make([]URL, 0, len(urlList))
	for _, urlStr := range urlList {
		res, err := url.Parse(urlStr)
		if err != nil {
			t.Fatalf("failed to parse URL: %s", err)
		}
		expectedURLs = append(expectedURLs, NewURL(res))
	}

	// unmarshal YAML
	var actualURLs []URL
	if err := yaml.Unmarshal([]byte(urlYAML), &actualURLs); err != nil {
		t.Fatalf("failed to unmarshal URLs: %s", err)
	}

	// Test
	assert.Equal(t, expectedURLs, actualURLs)
}