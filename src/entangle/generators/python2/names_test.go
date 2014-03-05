package python2

import (
	"testing"
)

func TestSnakeCaseString(t *testing.T) {
	for _, testCase := range []struct {
		Input string
		Expected string
	} {
		{"", ""},
		{"FooBar", "foo_bar"},
		{"HeadlineCNNNews", "headline_cnn_news"},
		{"CNN", "cnn"},
	} {
		actual := snakeCaseString(testCase.Input)
		if actual != testCase.Expected {
			t.Errorf("Expected snake casing of '%s' to return '%s', but it returned '%s'", testCase.Input, testCase.Expected, actual)
		}
	}
}
