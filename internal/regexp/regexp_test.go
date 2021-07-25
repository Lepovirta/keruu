package regexp

import (
	"regexp"
	"testing"

	"gopkg.in/yaml.v3"
	"github.com/stretchr/testify/assert"
)

const regexpYAML = `
- "^Sponsored Post:"
- \[ad\]$
`

var regexpList = []string{
	"^Sponsored Post:",
	"\\[ad\\]$",
}

func TestRegexpUnmarshallingToYAML(t *testing.T) {
	// setup expected regexps
	expectedRegexps := make([]RE, 0, len(regexpList))
	for _, reStr := range regexpList {
		re, err := regexp.Compile(reStr)
		if err != nil {
			t.Fatalf("failed to parse regexp: %s", err)
		}
		expectedRegexps = append(expectedRegexps, NewRE(re))
	}

	// unmarshal YAML
	var actualREs []RE
	if err := yaml.Unmarshal([]byte(regexpYAML), &actualREs); err != nil {
		t.Fatalf("failed to unmarshal regexps: %s", err)
	}

	// Test
	assert.Equal(t, expectedRegexps, actualREs)
}
