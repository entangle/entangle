// Package source provides a convience wrapper for reading source code for
// parsing and displaying error messages.
package source

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"unicode"
)

var (
	ErrUnicodeError = errors.New("Unicode error")
)

// Source.
type Source struct {
	// Data.
	data []rune

	// Lines.
	lines [][]rune

	// Path.
	path string
}

// Data of the source.
func (s *Source) Data() []rune {
	return s.data
}

// Line count of the source.
func (s *Source) LineCount() int {
	return len(s.lines)
}

// Line.
func (s *Source) Line(i int) string {
	return string(s.lines[i-1])
}

// Path.
func (s *Source) Path() string {
	return s.path
}

// Source from runes.
func FromRunes(data []rune, path string) (s *Source) {
	s = &Source{
		data:  data,
		lines: make([][]rune, 0),
		path:  path,
	}

	lineStart := 0
	for i, r := range data {
		if r == '\n' {
			l := data[lineStart:i]
			if len(l) > 0 && l[len(l)-1] == '\r' {
				l = l[:len(l)-1]
			}
			s.lines = append(s.lines, l)
			lineStart = i + 1
		}
	}

	s.lines = append(s.lines, data[lineStart:])

	return
}

// Source from reader.
func FromReader(reader io.Reader, path string) (s *Source, err error) {
	// Read all the runes from the input.
	var r rune
	var size int

	runes := make([]rune, 0, 10240)
	runeReader := bufio.NewReader(reader)

	for {
		if r, size, err = runeReader.ReadRune(); err != nil {
			if err == io.EOF {
				err = nil
				break
			}

			return
		}

		if r == unicode.ReplacementChar && size == 1 {
			err = ErrUnicodeError
			return
		}

		runes = append(runes, r)
	}

	return FromRunes(runes, path), nil
}

// Source from bytes.
func FromBytes(data []byte, path string) (s *Source, err error) {
	// Read all the runes from the input.
	var r rune
	var size int

	runes := make([]rune, 0, 10240)
	reader := bytes.NewBuffer(data)

	for {
		if r, size, err = reader.ReadRune(); err != nil {
			if err == io.EOF {
				err = nil
				break
			}

			return
		}

		if r == unicode.ReplacementChar && size == 1 {
			err = ErrUnicodeError
			return
		}

		runes = append(runes, r)
	}

	return FromRunes(runes, path), nil
}

// Source from string.
func FromString(data string, path string) (s *Source, err error) {
	return FromBytes([]byte(data), path)
}
