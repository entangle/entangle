package errors

import (
	"../token"
	"../source"
)

// Parse error frame.
type ParseErrorFrame struct {
	// Source.
	Source *source.Source

	// Start.
	Start token.Position

	// End.
	End token.Position
}

// Parse error.
type ParseError interface {
	error

	// Description.
	Description() string

	// Error frames.
	//
	// The last frame describes the actual error, while all previous frames
	// are guaranteed to describe imports.
	Frames() []ParseErrorFrame
}

// Parse error implementation.
type parseError struct {
	src *source.Source
	description string
	frames []ParseErrorFrame
}

func (p *parseError) Source() *source.Source {
	return p.src
}

func (p *parseError) Description() string {
	return p.description
}

func (p *parseError) Frames() []ParseErrorFrame {
	return p.frames
}

func (p *parseError) Error() string {
	return p.description
}

// New parse error for a token.
func NewParseErrorForToken(description string, tok *token.Token, src *source.Source, errorFrames []ParseErrorFrame) ParseError {
	err := &parseError {
		src: src,
		description: description,
		frames: make([]ParseErrorFrame, len(errorFrames) + 1),
	}

	copy(err.frames, errorFrames)
	err.frames[len(err.frames) - 1] = ParseErrorFrame {
		Source: src,
		Start: tok.Start,
		End: tok.End,
	}

	return err
}

// New parse error.
func NewParseError(description string, start token.Position, end token.Position, src *source.Source, errorFrames []ParseErrorFrame) ParseError {
	err := &parseError {
		src: src,
		description: description,
		frames: make([]ParseErrorFrame, len(errorFrames) + 1),
	}

	copy(err.frames, errorFrames)
	err.frames[len(err.frames) - 1] = ParseErrorFrame {
		Source: src,
		Start: start,
		End: end,
	}

	return err
}
