package lexer

import (
	"entangle/errors"
	"entangle/source"
	"entangle/token"
	"fmt"
	"testing"
)

type delimiterFixture struct {
	src       string
	tokenType token.TokenType
}

var delimiterFixtures = []delimiterFixture{
	{"", token.EndOfFile},
	{" ", token.EndOfFile},
	{"\r", token.EndOfFile},
	{"\t", token.EndOfFile},
	{"\n", token.NewLine},
	{"{", token.TokenType('{')},
	{"}", token.TokenType('}')},
	{"[", token.TokenType('[')},
	{"]", token.TokenType(']')},
	{"(", token.TokenType('(')},
	{")", token.TokenType(')')},
	{":", token.TokenType(':')},
	{".", token.TokenType('.')},
	{"/", token.TokenType('/')},
}

func assertLexerError(t *testing.T, stringSrc string, expected string) {
	src, err := source.FromString(stringSrc, "test.etg")
	if err != nil {
		panic(err)
	}

	lexer := NewLexer(src, []errors.ParseErrorFrame{})

	var tok token.Token

	for {
		tok, err = lexer.Lex()
		if err != nil {
			break
		}

		if tok.Type == token.EndOfFile {
			break
		}
	}

	if parseError, ok := err.(errors.ParseError); !ok || parseError.Description() != expected {
		t.Fatalf("Expected %v to occur, but error is %v", expected, err)
	}
}

// Test an expectedly valid source.
//
// Errors are logged from the function. In case of an error, tokens are nil.
func testValid(t *testing.T, stringSrc string) (tokens []token.Token) {
	// Create a source from the string.
	src, err := source.FromString(stringSrc, "<fixture>")
	if err != nil {
		t.Fatalf("source initialization failed for `%s`: %v", stringSrc, err)
	}

	// Create a lexer and parse out all the tokens.
	lexer := NewLexer(src, []errors.ParseErrorFrame{})

	var tok token.Token
	tokens = make([]token.Token, 0, 16)

	for {
		tok, err = lexer.Lex()
		if err != nil {
			break
		}

		tokens = append(tokens, tok)

		if tok.Type == token.EndOfFile {
			break
		}
	}

	// Assert that no errors occured.
	if err != nil {
		if parseErr, ok := err.(errors.ParseError); ok {
			lastFrame := parseErr.Frames()[len(parseErr.Frames())-1]

			t.Errorf("Unexpected error when parsing `%s`: %v, from %d:%d to %d:%d", stringSrc, err, lastFrame.Start.Line, lastFrame.Start.Character, lastFrame.End.Line, lastFrame.End.Character)
		} else {
			t.Errorf("Unexpected non parsing error when parsing `%s`: %v", stringSrc, err)
		}

		return nil
	}

	return tokens
}

// Test an expectedly valid source which should yield a first token of a given
// type.
//
// Behaves like testValid.
func testValidWithFirstOfType(t *testing.T, stringSrc string, tokenType token.TokenType) (tokens []token.Token) {
	tokens = testValid(t, stringSrc)
	if tokens == nil {
		return
	}

	if len(tokens) == 1 {
		t.Errorf("Expected lexing of `%s` to return at least one token other than EndOfFile", stringSrc)
		return
	}

	tok := &tokens[0]

	if tok.Type != tokenType {
		t.Errorf("Expected first lexed token from `%s` to be %s, but it's %s with value `%s`", stringSrc, tokenType, tok.Type, tok.StringValue)
		tokens = nil
	}

	return
}

func TestLexerLiterals(t *testing.T) {
	// * String literals.

	// Unexpected end of line returns errors.
	for _, prefix := range []string{
		"",
		"     ",
		"\n\n   \r\t",
	} {
		assertLexerError(t, fmt.Sprintf(`%s"horse`, prefix), errUnexpectedEndOfLineLiteralDesc)
		assertLexerError(t, fmt.Sprintf(`%s"horse\n`, prefix), errUnexpectedEndOfLineLiteralDesc)
		assertLexerError(t, fmt.Sprintf(`%s"horse\r\n`, prefix), errUnexpectedEndOfLineLiteralDesc)
	}
}
