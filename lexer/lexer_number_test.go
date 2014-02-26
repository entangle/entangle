package lexer

import (
	"../token"
	"fmt"
	"testing"
)

func testValidUintConstant(t *testing.T, stringSrc string, value uint64) {
	for _, delFixture := range delimiterFixtures {
		if delFixture.src == "." {
			continue
		}

		fixtureSrc := fmt.Sprintf("%s%s", stringSrc, delFixture.src)

		tokens := testValidWithFirstOfType(t, fixtureSrc, token.UintConstant)
		if tokens == nil {
			continue
		}

		tok := &tokens[0]
		if tok.UintValue != value {
			t.Fatalf("Expected token unsigned integer value lexed from `%s` to be %d, but it is %d", fixtureSrc, value, tok.UintValue)
		}

		if tokens[1].Type != delFixture.tokenType {
			t.Errorf("Expected token following unsigned integer constant when lexing `%s` to be of type %s but it is of type %s", fixtureSrc, delFixture.tokenType, tokens[1].Type)
		}
	}
}

func testValidIntConstant(t *testing.T, stringSrc string, value int64) {
	for _, delFixture := range delimiterFixtures {
		if delFixture.src == "." {
			continue
		}

		fixtureSrc := fmt.Sprintf("%s%s", stringSrc, delFixture.src)

		tokens := testValidWithFirstOfType(t, fixtureSrc, token.IntConstant)
		if tokens == nil {
			continue
		}

		tok := &tokens[0]
		if tok.IntValue != value {
			t.Fatalf("Expected token integer value lexed from `%s` to be %d, but it is %d", fixtureSrc, value, tok.IntValue)
		}

		if tokens[1].Type != delFixture.tokenType {
			t.Errorf("Expected token following integer constant when lexing `%s` to be of type %s but it is of type %s", fixtureSrc, delFixture.tokenType, tokens[1].Type)
		}
	}
}

func testValidFloatConstant(t *testing.T, stringSrc string, value float64) {
	for _, delFixture := range delimiterFixtures {
		if delFixture.src == "." {
			continue
		}

		fixtureSrc := fmt.Sprintf("%s%s", stringSrc, delFixture.src)

		tokens := testValidWithFirstOfType(t, fixtureSrc, token.FloatConstant)
		if tokens == nil {
			continue
		}

		tok := &tokens[0]
		if tok.FloatValue != value {
			t.Fatalf("Expected token floating point value lexed from `%s` to be %f, but it is %f", fixtureSrc, value, tok.FloatValue)
		}

		if tokens[1].Type != delFixture.tokenType {
			t.Errorf("Expected token following floating point constant when lexing `%s` to be of type %s but it is of type %s", fixtureSrc, delFixture.tokenType, tokens[1].Type)
		}
	}
}

func TestLexerHexadecimalConstant(t *testing.T) {
	testValidUintConstant(t, "0x1", 1)
	testValidUintConstant(t, "0x0123456789abcdef", 0x0123456789abcdef)
	testValidIntConstant(t, "-0x0123456789abcdef", -0x0123456789abcdef)
}

func TestLexerOctalConstant(t *testing.T) {
	testValidUintConstant(t, "01", 1)
	testValidUintConstant(t, "01234567", 01234567)
	testValidIntConstant(t, "-01234567", -01234567)
}

func TestLexerDecimalConstant(t *testing.T) {
	testValidUintConstant(t, "1", 1)
	testValidUintConstant(t, "123456789", 123456789)
	testValidIntConstant(t, "-123456789", -123456789)
}

func TestLexerFloatConstant(t *testing.T) {
	// * Floating point numbers.
	testValidFloatConstant(t, "0.", 0.0)
	testValidFloatConstant(t, ".0", 0.0)
	testValidFloatConstant(t, "1.0", 1.0)
	testValidFloatConstant(t, "1.0e5", 1.0e5)
	testValidFloatConstant(t, "1.797693134862315708145274237317043567981e308", 1.797693134862315708145274237317043567981e308)
	testValidFloatConstant(t, "4.940656458412465441765687928682213723651e-324", 4.940656458412465441765687928682213723651e-324)
	testValidFloatConstant(t, "-00123456789.0123456789e123", -123456789.0123456789e123)
	testValidFloatConstant(t, "00123456789.0123456789e123", 123456789.0123456789e123)
}
