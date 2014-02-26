package lexer

import (
	"../token"
	"fmt"
	"testing"
)

type identifierFixture struct {
	identifier string
	tokenType  token.TokenType
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
		for _, delFixture := range delimiterFixtures {
			stringSrc := fmt.Sprintf("%s%s", fixture.identifier, delFixture.src)
			if tokens := testValidWithFirstOfType(t, stringSrc, fixture.tokenType); tokens != nil {
				if tokens[0].StringValue != fixture.identifier {
					t.Errorf("Expected first token when lexing `%s` to have a string value of `%s`, but it has a value of `%s`", stringSrc, fixture.identifier, tokens[0].StringValue)
				}

				if tokens[1].Type != delFixture.tokenType {
					t.Errorf("Expected token following keyword when lexing `%s` to be of type %s but it is of type %s", stringSrc, delFixture.tokenType, tokens[1].Type)
				}
			}
		}

		// Invalid character following identifier or keyword returns error.
		for _, suffix := range []string{
			"\"",
		} {
			assertLexerError(t, fmt.Sprintf("%s%s", fixture.identifier, suffix), errUnexpectedCharacterDesc)
		}
	}
}
