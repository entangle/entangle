package utils

import (
	"regexp"
	"strings"
)

const (
	WHITESPACE = "\t\n\x0b\x0c\r "
)

// Split a string using a regular expression.
func regexpSplit(expression *regexp.Regexp, text string) []string {
	indexes := expression.FindAllStringIndex(text, -1)
	laststart := 0
	result := make([]string, 0, len(indexes)*2+1)
	for _, element := range indexes {
		if laststart != element[0] {
			result = append(result, text[laststart:element[0]])
		}
		if element[0] != element[1] {
			result = append(result, text[element[0]:element[1]])
		}
		laststart = element[1]
	}
	if laststart < len(text) {
		result = append(result, text[laststart:])
	}
	return result
}

// Text wrapper.

// Configuration and control structure for wrapping and filling text.
type TextWrapper struct {
	// Text width.

	// Maximum line width.
	Width int

	// Expand tabs.

	// Whether or not to expand tabs into spaces before processing. A tab will
	// become [1 ; 8] spaces depending on its position in the line when.
	// expanded. If disabled, a tab will be treated as one character and may
	// produce unexpected results.
	ExpandTabs bool

	// Normalize whitespace.

	// Whether or not to replace all whitespaces with a space character (` `.)
	NormalizeWhitespace bool

	// Break overflowing words.

	// When enabled, words that will overflow a single line is broken.
	BreakOverflowingWords bool

	// Break on hyphnes.
	BreakOnHyphens bool

	// Trim whitespace.

	// Whether or not to trim whitespace from resulting lines.
	TrimWhitespace bool
}

// Default text wrapper.
var DefaultTextWrapper = TextWrapper{
	Width:                 79,
	ExpandTabs:            true,
	NormalizeWhitespace:   true,
	BreakOverflowingWords: true,
	BreakOnHyphens:        true,
	TrimWhitespace:        true,
}

var chunkSeparatorPattern = regexp.MustCompile(`((?:--+)|\s+)`)
var hyphenatedWordPattern = regexp.MustCompile(`^(\w+-)\w`)
var simpleChunkSeparatorPattern = regexp.MustCompile(`(\s+)`)

func NewSimpleTextWrapper(width int) *TextWrapper {
	return &TextWrapper{
		Width:                 width,
		ExpandTabs:            true,
		NormalizeWhitespace:   true,
		BreakOverflowingWords: true,
		BreakOnHyphens:        true,
		TrimWhitespace:        true,
	}
}

func (t *TextWrapper) Wrap(input string) (lines []string) {
	// Transform whitespaces if necessary.
	if t.ExpandTabs {
		input = ExpandTabs(input, 8)
	}
	if t.NormalizeWhitespace {
		input = strings.Map(func(in rune) rune {
			if strings.ContainsRune(WHITESPACE, in) {
				return ' '
			}
			return in
		}, input)
	}

	// Split the text into chunks we can deal with.
	var chunks []string
	if t.BreakOnHyphens {
		// This is not particularly beautiful nor efficient, but the lack
		// of look-ahead and -behind in the RE2 specification causes trouble
		// for us here. This should work, though ;)
		intermediateChunks := regexpSplit(simpleChunkSeparatorPattern, input)
		chunks = make([]string, 0, len(intermediateChunks))

		for _, c := range intermediateChunks {
			if !hyphenatedWordPattern.MatchString(c) {
				chunks = append(chunks, c)
				continue
			}

			hyphenIndex := strings.IndexRune(c, '-')
			for hyphenIndex != -1 && hyphenIndex != len(c)-1 {
				chunks = append(chunks, c[0:hyphenIndex+1])
				c = c[hyphenIndex+1:]
				hyphenIndex = strings.IndexRune(c, '-')
			}

			if len(c) > 0 {
				chunks = append(chunks, c)
			}
		}
	} else {
		chunks = regexpSplit(simpleChunkSeparatorPattern, input)
	}

	// Wrap the chunks.
	lines = []string{}

	i := 0
	for i < len(chunks) {
		// Initialize the line state.
		lineChunks := []string{}
		lineLength := 0

		// Drop initial whitespace chunks if necessary.
		if t.TrimWhitespace {
			for i < len(chunks) {
				if len(strings.TrimSpace(chunks[i])) != 0 {
					break
				}
				i++
			}
		}

		// Iterate across the residual chunks and try to fit them on the line.
		for i < len(chunks) {
			c := chunks[i]
			l := len(c)
			if lineLength+l > t.Width {
				break
			}

			lineChunks = append(lineChunks, c)
			lineLength += l
			i++
		}

		// If the current line is full, and the next chunk is too big to fit on
		// any line, handle this case.
		if i < len(chunks) && len(chunks[i]) > t.Width {
			if t.BreakOverflowingWords {
				// If we break overflowing words, try to take a bit of the
				// chunk this time around.
				residualLength := t.Width - lineLength
				if residualLength > 0 {
					lineChunks = append(lineChunks, chunks[i][0:residualLength])
					chunks[i] = chunks[i][residualLength:]
				}
			} else if lineLength == 0 {
				// Take the whole chunk if this is a fresh line.
				lineChunks = append(lineChunks, chunks[i])
				i++
			}
		}

		// Remove trailing white space chunks if necessary.
		if t.TrimWhitespace {
			for len(lineChunks) > 0 {
				if len(strings.TrimSpace(lineChunks[len(lineChunks)-1])) != 0 {
					break
				}
				lineChunks = lineChunks[0 : len(lineChunks)-1]
			}
		}

		// Append the line if there's anything in it.
		if len(lineChunks) > 0 {
			lines = append(lines, strings.Join(lineChunks, ""))
		}
	}

	return
}
