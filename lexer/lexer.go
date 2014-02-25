// Package lexer provides the lexer for parsing Entangle IDL files.
package lexer

import (
	"../errors"
	"../source"
	"../token"
	"fmt"
	"strconv"
)

var (
	errExpectedHexadecimalDigitDesc   = "expected hexadecimal digit"
	errUnexpectedCharacterDesc        = "unexpected character"
	errUnexpectedEndOfLineLiteralDesc = "unexpected end of line in literal"
	errUnexpectedEndOfLineNumberDesc  = "unexpected end of line in numerical"
	errNumberOutOfRangeDesc           = "number is out of range"
)

const (
	eof = rune(-1)
	eol = rune('\n')
)

// Lexer.
type Lexer struct {
	// Source.
	src *source.Source

	// Data.
	data []rune

	// Data position.
	dataPosition int

	// Data length.
	dataLength int

	// Whether there has previously been a token present on the line.
	lineHasHadToken bool

	// Line index.
	//
	// 1-based.
	lineIndex int

	// Current line position.
	//
	// 0-based. A negative value indicates that no line has been read yet.
	linePosition int

	// Previous character.
	//
	// NUL if the current character is the first character.
	prev rune

	// Current character.
	cur rune

	// Current position.
	position token.Position

	// Error frames.
	errorFrames []errors.ParseErrorFrame
}

// Source.
func (l *Lexer) Source() *source.Source {
	return l.src
}

// Parse error.
func (l *Lexer) parseError(description string, start token.Position, end token.Position) error {
	return errors.NewParseError(description, start, end, l.src, l.errorFrames)
}

// Parse error at the current position.
func (l *Lexer) parseErrorHere(description string) error {
	return errors.NewParseError(description, l.position, l.position, l.src, l.errorFrames)
}

// Read the next character from the scanner.
func (l *Lexer) next() {
	if l.dataPosition == l.dataLength {
		l.cur = eof
		return
	}

	if l.cur == eol {
		l.position.Line++
		l.position.Character = 1
		l.lineHasHadToken = false
	} else {
		l.position.Character++
	}

	l.prev = l.cur
	l.cur = l.data[l.dataPosition]
	l.dataPosition++

	return
}

// Peek ahead to the next character.
func (l *Lexer) peek() rune {
	if l.dataPosition == l.dataLength {
		return eof
	}

	return l.data[l.dataPosition]
}

// Test if a rune is a white space character.
func isWhitespace(r rune) bool {
	return r != eof && r < 128 && whitespaceCharacterTable[r]
}

// Test if a rune is a digit.
func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

// Test if a rune is a hexadecimal digit.
func isHexadecimalDigit(r rune) bool {
	return r > 0 && r < 128 && hexDigitCharacterTable[r]
}

// Test if a rune is a valid delimiter following an identifier or a numerical
// value.
func isValidDelimiter(r rune) bool {
	return r == eof || r < 128 && delimiterCharacterTable[r]
}

// Skip whitespace.
func (l *Lexer) skipWhitespace() {
	for isWhitespace(l.cur) {
		l.next()
	}
}

// Skip a multi line comment.
func (l *Lexer) skipMultiLineComment() {
	for {
		l.next()

		switch l.cur {
		case eof:
			return

		case '/':
			if l.prev == '*' {
				l.next()
				return
			}
		}
	}
}

// Parse a single line comment.
func (l *Lexer) parseSingleLineComment() (t token.Token) {
	t.Start = l.position
	t.Type = token.DocumentationLine
	dataStart := l.dataPosition

	for {
		l.next()

		switch l.cur {
		case eof, eol:
			t.StringValue = string(l.data[dataStart:l.dataPosition])
			t.End = l.position.Before()
			return
		}
	}
}

// Skip a single line comment.
func (l *Lexer) skipSingleLineComment() {
	for {
		l.next()

		switch l.cur {
		case eof, eol:
			return
		}
	}
}

// Error indicating that the current character was unexpected.
func (l *Lexer) unexpectedCharacter() (t token.Token, err error) {
	err = l.parseErrorHere(errUnexpectedCharacterDesc)
	return
}

// Parse an integer to a token.
func (l *Lexer) parseIntToToken(t *token.Token, negative bool, rawValue string, base int) (err error) {
	if negative {
		var intValue int64
		t.Type = token.IntConstant
		intValue, err = strconv.ParseInt(rawValue, base, 64)
		t.IntValue = -intValue
	} else {
		t.Type = token.UintConstant
		t.UintValue, err = strconv.ParseUint(rawValue, base, 64)
	}

	if err != nil {
		if err != strconv.ErrRange {
			panic(fmt.Sprintf("Unxpected error: %s", err))
		}

		err = l.parseError(errNumberOutOfRangeDesc, t.Start, t.End)
	}

	return
}

// Parse a number.
//
// Invoked with a valid starting character for a number.
func (l *Lexer) parseNumber() (t token.Token, err error) {
	t.Start = l.position

	// Parse the sign if present.
	negative := false
	dataStart := l.dataPosition - 1
	dataMantissaStart := dataStart

	if l.cur == '+' || l.cur == '-' {
		negative = l.cur == '-'
		dataMantissaStart++
		l.next()
	}

	// Let's check immdiately if we're parsing a hexadecimal number by peeking
	// ahead.
	if next := l.peek(); l.cur == '0' && (next == 'x' || next == 'X') {
		// Skip past the '0' and the 'x'.
		l.next()
		l.next()

		// Read as many valid characters as possible.
		for isHexadecimalDigit(l.cur) {
			l.next()
		}

		if !isValidDelimiter(l.cur) {
			return l.unexpectedCharacter()
		}

		// Attempt to convert the number.
		rawValue := string(l.data[dataMantissaStart+2 : l.dataPosition-1])

		if len(rawValue) == 0 {
			if l.cur == eof || l.cur == eol {
				err = l.parseErrorHere(errUnexpectedEndOfLineNumberDesc)
			} else {
				err = l.parseErrorHere(errExpectedHexadecimalDigitDesc)
			}

			return
		}

		t.StringValue = string(l.data[dataStart : l.dataPosition-1])
		t.End = l.position.Before()

		err = l.parseIntToToken(&t, negative, rawValue, 16)
		return
	}

	// Let's assume that we're parsing a mantissa and then let's reduce to
	// a more specific case as soon as possible. Basically, we can only
	// determine if we're parsing either a decimal integer or an octal integer
	// when we know that we're not going to hit a period ('.') or an exponent
	// prefix ('e' or 'E'.) Until then, all bets are on.
	digitsBeforeDot := 0

	for isDigit(l.cur) {
		digitsBeforeDot++
		l.next()
	}

	// If the next character is not a dot or an exponent prefix, we're dealing
	// with an octal or decimal integer.
	if l.cur != '.' && l.cur != 'e' && l.cur != 'E' {
		// Make sure that we have any digits at all.
		if digitsBeforeDot == 0 {
			t, err = l.unexpectedCharacter()

			if l.cur == eol || l.cur == eof {
				err = l.parseErrorHere(errUnexpectedEndOfLineNumberDesc)
			}

			return
		}

		// Make sure that we've hit a valid delimiter.
		if !isValidDelimiter(l.cur) {
			return l.unexpectedCharacter()
		}

		// Parse as an octal or decimal number.
		rawValue := string(l.data[dataMantissaStart : l.dataPosition-1])
		t.StringValue = string(l.data[dataStart : l.dataPosition-1])
		t.End = l.position.Before()

		base := 10
		if rawValue[0] == '0' {
			base = 8

			// In the case of an octal, we need to verify that all digits are
			// valid octal digits or present an error in the case it's not.
			for i, c := range rawValue {
				if c > '7' {
					pos := t.Start
					pos.Character += dataStart - dataMantissaStart + i
					err = l.parseError(errUnexpectedCharacterDesc, pos, pos)
					return
				}
			}
		}

		err = l.parseIntToToken(&t, negative, rawValue, base)
		return
	}

	// Given that we've gotten here, we're dealing with a floating point,
	// number, so let's parse any digits after the dot if we've reached a dot.
	digitsAfterDot := 0

	if l.cur == '.' {
		l.next()

		for isDigit(l.cur) {
			digitsAfterDot++
			l.next()
		}
	}

	// Check that we have a valid mantissa at this point.
	if digitsBeforeDot == 0 && digitsAfterDot == 0 {
		t, err = l.unexpectedCharacter()

		if l.cur == eol || l.cur == eof {
			err = l.parseErrorHere(errUnexpectedEndOfLineNumberDesc)
		}

		return
	}

	// Read the exponent if available.
	if l.cur == 'e' || l.cur == 'E' {
		l.next()
		exponentDigits := 0

		if l.cur == '-' || l.cur == '+' {
			l.next()
		}

		for isDigit(l.cur) {
			exponentDigits++
			l.next()
		}

		if exponentDigits == 0 {
			t, err = l.unexpectedCharacter()

			if l.cur == eol || l.cur == eof {
				err = l.parseErrorHere(errUnexpectedEndOfLineNumberDesc)
			}

			return
		}
	}

	// Make sure that we've hit a valid delimiter at this point.
	if !isValidDelimiter(l.cur) {
		return l.unexpectedCharacter()
	}

	// Attempt to parse the floating point value.
	t.StringValue = string(l.data[dataStart : l.dataPosition-1])
	t.End = l.position.Before()
	t.Type = token.FloatConstant

	if t.FloatValue, err = strconv.ParseFloat(t.StringValue, 64); err != nil {
		if err != strconv.ErrRange {
			panic(fmt.Sprintf("Unexpected error: %s", err))
		}

		err = l.parseError(errNumberOutOfRangeDesc, t.Start, t.End)
	}

	return
}

// Parse a quoted string.
func (l *Lexer) parseQuotedString() (t token.Token, err error) {
	escaped := false
	chars := make([]rune, 0, 128)
	t.Start = l.position
	t.Type = token.Literal

	for {
		cur := l.cur

		if escaped {
			switch cur {
			case eof, eol:
				err = l.parseError(errUnexpectedEndOfLineLiteralDesc, t.Start, l.position)
				return

			case 'n':
				chars = append(chars, eol)

			case 'r':
				chars = append(chars, '\r')

			case 't':
				chars = append(chars, '\t')

			default:
				chars = append(chars, cur)
			}

			escaped = false
		} else {
			switch cur {
			case '\\':
				escaped = true

			case '"':
				t.End = l.position
				t.StringValue = string(chars)
				l.next()
				return

			case eof, eol:
				err = l.parseError(errUnexpectedEndOfLineLiteralDesc, t.Start, l.position)
				return

			default:
				chars = append(chars, cur)
			}
		}

		l.next()
	}
}

// Lex.
//
// In case of an error, token is populated with information about the
// problematic data.
func (l *Lexer) Lex() (t token.Token, err error) {
	for {
		// Skip whitespace.
		l.skipWhitespace()
		firstTokenInLine := !l.lineHasHadToken
		l.lineHasHadToken = true

		// Handle single character tokens.
		t.Start = l.position
		t.End = l.position
		t.Type = token.TokenType(l.cur)
		t.StringValue = string(l.cur)

		switch current := l.cur; {
		case current == eof:
			t.Type = token.EndOfFile
			return

		case current < 128 && reservedCharacterTable[uint8(current)]:
			// Move on to the next character.
			l.next()

			// Handle reserved characters.
			switch current {
			case eol:
				t.Type = token.NewLine

			case '"':
				return l.parseQuotedString()

			case '/':
				switch l.cur {
				case '/':
					if firstTokenInLine {
						return l.parseSingleLineComment(), nil
					} else {
						l.skipSingleLineComment()
						continue
					}

				case '*':
					l.skipMultiLineComment()
					continue
				}
			}

			return

		case current < 128 && identifierStartCharacterTable[uint8(current)]:
			// Parse the identifier.
			return l.parseIdentifier()

		case current < 128 && numericalFirstCharacterTable[uint8(current)]:
			// If the first character is a dot ('.') and the next character
			// is not valid for a digit, let's return the dot as a special
			// character.
			if current == '.' && !isDigit(l.peek()) {
				l.next()
				return
			}

			return l.parseNumber()

		default:
			// Move on to the next charcter.
			l.next()

			// Pass through the raw character.
			return
		}
	}

	return
}

// Parse an identifier.
//
// Invoked with the guarantee that the current character is a valid start
// character for an identifier.
func (l *Lexer) parseIdentifier() (t token.Token, err error) {
	runes := append(make([]rune, 0, 128), l.cur)
	t.Start = l.position

	for {
		l.next()

		if l.cur == eof || l.cur > 127 || !identifierCharacterTable[l.cur] {
			break
		}

		runes = append(runes, l.cur)
	}

	// Make sure that an identifier is followed by either a whitespace, a
	// new line, and end of line or a control character.
	if !isValidDelimiter(l.cur) {
		return l.unexpectedCharacter()
	}

	// Construct a string and determine if it's a keyword or "just" an
	// identifier.
	t.End = l.position.Before()
	identifier := string(runes)

	t.Type = IdentifierTokenType(identifier)
	t.StringValue = identifier

	return
}

// New lexer.
//
// The provided error frames are used for reporting errors.
func NewLexer(src *source.Source, errorFrames []errors.ParseErrorFrame) (l *Lexer) {
	// Create the lexer.
	l = &Lexer{
		src:        src,
		data:       src.Data(),
		dataLength: len(src.Data()),
		prev:       0,
		cur:        ' ',
		position: token.Position{
			Line:      1,
			Character: 0,
		},
		errorFrames: errorFrames,
	}

	return
}
