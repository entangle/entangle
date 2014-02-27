package utils

import (
	"strings"
)

// Expand tabs.

// @param input Text to have tabs expanded.
// @param tabWidth Tab width.
func ExpandTabs(input string, tabWidth int) (result string) {
	result = ""
	sliceStart := 0

	for i, r := range input {
		if r != '\t' {
			continue
		}

		if sliceStart < i {
			result += input[sliceStart:i]
		}
		result += strings.Repeat(" ", tabWidth-(len(result)%tabWidth))
		sliceStart = i + 1
	}

	if sliceStart < len(input) {
		result += input[sliceStart:]
	}

	return
}

// Mask all non-whitespace with whitespace and expand tabs.
func MaskWithWhitespaceExpanded(input string, tabWidth int) (result string) {
	sliceStart := 0
	spaces := 0

	for i, r := range input {
		if r != '\t' {
			continue
		}

		if sliceStart < i {
			spaces += i - sliceStart
		}
		spaces += tabWidth - (spaces % tabWidth)
		sliceStart = i + 1
	}

	if sliceStart < len(input) {
		spaces += len(input) - sliceStart
	}

	return strings.Repeat(" ", spaces)
}
