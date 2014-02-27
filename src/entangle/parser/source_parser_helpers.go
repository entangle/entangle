package parser

import (
	"entangle/token"
	"strings"
	"unicode"
)

// Expect a rune.
//
// If the rune is found, the source is moved ahead one token as the rune has
// been validated. Context description will be used as
// "... in <context description>".
func (p *sourceParser) expectRune(r rune, contextDesc string) (err error) {
	switch p.tok.Type {
	case token.TokenType(r):
		return p.next()

	case token.NewLine:
		return p.parseErrorHeref("unexpected new line in %s, expected '%s'", contextDesc, r)

	case token.EndOfFile:
		return p.parseErrorHeref("unexpected end of file in %s, expected '%s'", contextDesc, r)

	default:
		return p.parseErrorHeref("expected '%s' in %s", string(r), contextDesc)
	}
}

// Skip new lines.
func (p *sourceParser) skipNewLines() (err error) {
	for p.tok.Type == token.NewLine {
		if err = p.next(); err != nil {
			return
		}
	}

	return
}

// Move to the next token and skip new lines from there.
func (p *sourceParser) nextAndSkipNewLines() (err error) {
	if err = p.next(); err != nil {
		return
	}

	return p.skipNewLines()
}

// Skip new lines and store documentation lines as we go along.
func (p *sourceParser) skipNewLinesStoreDocumentation() (err error) {
	p.documentationLines = []token.Token{}

	for {
		if p.tok.Type == token.NewLine {
			// Reset the documentation cache if the previous token was not
			// a documentation line.
			if p.prev.Type != token.DocumentationLine {
				p.documentationLines = []token.Token{}
			}
		} else if p.tok.Type == token.DocumentationLine {
			// Store the documentation line.
			p.documentationLines = append(p.documentationLines, p.tok)
		} else {
			break
		}

		if err = p.next(); err != nil {
			return
		}
	}

	return
}

// Get the documentation paragraphs.
//
// Reading the documentation paragraphs clears the stored documentation lines
// thus making it easy to use.
func (p *sourceParser) documentationParagraphs() (paragraphs []string) {
	paragraphs = make([]string, 0, len(p.documentationLines))

	if len(p.documentationLines) == 0 {
		return
	}

	paragraphSegments := []string{}

	for _, t := range p.documentationLines {
		segment := t.StringValue
		trimmedSegment := strings.TrimSpace(segment)

		if len(trimmedSegment) == 0 {
			if len(paragraphSegments) > 0 {
				paragraphs = append(paragraphs, strings.Join(paragraphSegments, " "))
				paragraphSegments = []string{}
			}
		} else {
			if segment[0] == ' ' {
				segment = segment[1:]
			}

			segment = strings.TrimRightFunc(segment, unicode.IsSpace)

			paragraphSegments = append(paragraphSegments, segment)
		}
	}

	if len(paragraphSegments) > 0 {
		paragraphs = append(paragraphs, strings.Join(paragraphSegments, " "))
	}

	p.documentationLines = []token.Token{}

	return
}
