package utils

import (
	"testing"
)

func testSimpleTextWrapperWrapCase(t *testing.T, input string, width int, expected []string) {
	actual := NewSimpleTextWrapper(width).Wrap(input)
	if len(actual) != len(expected) {
		t.Fatalf("expected %d lines of wrapped text, but got %d: %s", len(expected), len(actual), actual)
	}

	for i, actualLine := range actual {
		if actualLine != expected[i] {
			t.Fatalf("line %d was expected to be '%s' but is '%s'", i+1, expected[i], actualLine)
		}
	}
}

func TestSimpleTextWrapperWrap(t *testing.T) {
	text := "Hello there, how are you this fine day?  I'm glad to hear it!"

	testSimpleTextWrapperWrapCase(t, text, 12, []string{
		"Hello there,",
		"how are you",
		"this fine",
		"day?  I'm",
		"glad to hear",
		"it!",
	})

	testSimpleTextWrapperWrapCase(t, text, 42, []string{
		"Hello there, how are you this fine day?",
		"I'm glad to hear it!",
	})
}
