package lexer

import (
	"../errors"
	"../source"
	"../token"
	"fmt"
	"testing"
)

type identifierFixture struct {
	identifier string
	tokenType  token.TokenType
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

func assertLexerSuccess(t *testing.T, stringSrc string) (tokens []token.Token) {
	src, err := source.FromString(stringSrc, "test.etg")
	if err != nil {
		panic(err)
	}

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

	if err != nil {
		t.Fatalf("Unexpected error: %v, at %d:%d", err, tok.Start.Line, tok.Start.Character)
	}

	return tokens
}

func assertTokenType(t *testing.T, tok token.Token, tokenType token.TokenType) {
	if tok.Type != tokenType {
		t.Fatalf("Expected token to be of type %s but it is of type %s", tokenType, tok.Type)
	}
}

func assertIntConstantToken(t *testing.T, tok token.Token, value int64) {
	assertTokenType(t, tok, token.IntConstant)
	if tok.IntValue != value {
		t.Fatalf("Expected token integer value to be %d, but it is %d", value, tok.IntValue)
	}
}

func assertUintConstantToken(t *testing.T, tok token.Token, value uint64) {
	assertTokenType(t, tok, token.UintConstant)
	if tok.UintValue != value {
		t.Fatalf("Expected token unsigned integer value to be %u, but it is %u", value, tok.UintValue)
	}
}

func assertFloatConstantToken(t *testing.T, tok token.Token, value float64) {
	assertTokenType(t, tok, token.FloatConstant)
	if tok.FloatValue != value {
		t.Fatalf("Expected token floating point value to be %u, but it is %u", value, tok.FloatValue)
	}
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

func TestLexerIdentifiers(t *testing.T) {
	// * Identifiers and keywords.
	for _, fixture := range []identifierFixture{
		{"horse", token.Identifier},
		{"Horse", token.Identifier},
		{"A0", token.Identifier},
		{"HORSE_STUFF", token.Identifier},
		{"int64", token.Int64},
		{"uint64", token.Uint64},
		{"float64", token.Float64},
		{"int32", token.Int32},
		{"uint32", token.Uint32},
		{"float32", token.Float32},
		{"import", token.Import},
		{"typedef", token.Typedef},
		{"int8", token.Int8},
		{"uint8", token.Uint8},
		{"int16", token.Int16},
		{"uint16", token.Uint16},
		{"struct", token.Struct},
		{"service", token.Service},
		{"enum", token.Enum},
		{"binary", token.Binary},
		{"bool", token.Bool},
		{"const", token.Const},
	} {
		// Valid character or end of file following identifier or keyword
		// succeeds.
		for _, suffix := range []string{
			"",
			" ",
			"\r",
			"\t",
			"\n",
			"{",
			"}",
			"[",
			"]",
			"(",
			")",
			":",
			".",
			"/",
		} {
			t.Logf("Validating that parsing of `%s%s` succeeds", fixture.identifier, suffix)
			assertLexerSuccess(t, fmt.Sprintf("%s%s", fixture.identifier, suffix))
		}

		// Invalid character following identifier or keyword returns error.
		for _, suffix := range []string{
			"\"",
		} {
			t.Logf("Validating that parsing of `%s%s` fails with ErrUnexpectedCharacter", fixture.identifier, suffix)
			assertLexerError(t, fmt.Sprintf("%s%s", fixture.identifier, suffix), errUnexpectedCharacterDesc)
		}
	}
}

func TestLexerNumbers(t *testing.T) {
	// * Hexadecimal integers.
	tokens := assertLexerSuccess(t, "0x1")
	assertUintConstantToken(t, tokens[0], 1)

	tokens = assertLexerSuccess(t, "0x0123456789abcdef")
	assertUintConstantToken(t, tokens[0], 0x0123456789abcdef)

	tokens = assertLexerSuccess(t, "-0x0123456789abcdef")
	assertIntConstantToken(t, tokens[0], -0x0123456789abcdef)

	// * Octal integers.
	tokens = assertLexerSuccess(t, "01")
	assertUintConstantToken(t, tokens[0], 1)

	tokens = assertLexerSuccess(t, "01234567")
	assertUintConstantToken(t, tokens[0], 01234567)

	tokens = assertLexerSuccess(t, "-01234567")
	assertIntConstantToken(t, tokens[0], -01234567)

	// * Decimal integers.
	tokens = assertLexerSuccess(t, "1")
	assertUintConstantToken(t, tokens[0], 1)

	tokens = assertLexerSuccess(t, "123456789")
	assertUintConstantToken(t, tokens[0], 123456789)

	tokens = assertLexerSuccess(t, "-123456789")
	assertIntConstantToken(t, tokens[0], -123456789)

	// * Floating point numbers.
	tokens = assertLexerSuccess(t, "0.")
	assertFloatConstantToken(t, tokens[0], 0.0)

	tokens = assertLexerSuccess(t, ".0")
	assertFloatConstantToken(t, tokens[0], 0.0)

	tokens = assertLexerSuccess(t, "1.0")
	assertFloatConstantToken(t, tokens[0], 1.0)

	tokens = assertLexerSuccess(t, "1.0e5")
	assertFloatConstantToken(t, tokens[0], 1.0e5)

	tokens = assertLexerSuccess(t, "1.797693134862315708145274237317043567981e308")
	assertFloatConstantToken(t, tokens[0], 1.797693134862315708145274237317043567981e308)

	tokens = assertLexerSuccess(t, "4.940656458412465441765687928682213723651e-324")
	assertFloatConstantToken(t, tokens[0], 4.940656458412465441765687928682213723651e-324)

	tokens = assertLexerSuccess(t, "-00123456789.0123456789e123")
	assertFloatConstantToken(t, tokens[0], -123456789.0123456789e123)

	tokens = assertLexerSuccess(t, "00123456789.0123456789e123")
	assertFloatConstantToken(t, tokens[0], 123456789.0123456789e123)
}
