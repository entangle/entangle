package lexer

import (
	"../token"
	"testing"
)

func assertValidIdentifierTokenType(t *testing.T, input string, tokenType token.TokenType) {
	result := IdentifierTokenType(input)
	if result != tokenType {
		t.Fatalf("unexpected keyword token for %s: %d", input, tokenType)
	}
}

func TestIdentifierTokenType(t *testing.T) {
	// Check that all valid keywords return the expected token.
	assertValidIdentifierTokenType(t, "import", token.Import)
	assertValidIdentifierTokenType(t, "bool", token.Bool)
	assertValidIdentifierTokenType(t, "binary", token.Binary)
	assertValidIdentifierTokenType(t, "float32", token.Float32)
	assertValidIdentifierTokenType(t, "float64", token.Float64)
	assertValidIdentifierTokenType(t, "int8", token.Int8)
	assertValidIdentifierTokenType(t, "int16", token.Int16)
	assertValidIdentifierTokenType(t, "int32", token.Int32)
	assertValidIdentifierTokenType(t, "int64", token.Int64)
	assertValidIdentifierTokenType(t, "uint8", token.Uint8)
	assertValidIdentifierTokenType(t, "uint16", token.Uint16)
	assertValidIdentifierTokenType(t, "uint32", token.Uint32)
	assertValidIdentifierTokenType(t, "uint64", token.Uint64)
	assertValidIdentifierTokenType(t, "const", token.Const)
	assertValidIdentifierTokenType(t, "enum", token.Enum)
	assertValidIdentifierTokenType(t, "struct", token.Struct)
	assertValidIdentifierTokenType(t, "typedef", token.Typedef)
	assertValidIdentifierTokenType(t, "service", token.Service)
}
