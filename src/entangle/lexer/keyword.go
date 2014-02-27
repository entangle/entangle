package lexer

import (
	"entangle/token"
)

var keywordMap = map[string]token.TokenType{
	"int64":     token.Int64,
	"uint64":    token.Uint64,
	"float64":   token.Float64,
	"int32":     token.Int32,
	"uint32":    token.Uint32,
	"float32":   token.Float32,
	"import":    token.Import,
	"typedef":   token.Typedef,
	"int8":      token.Int8,
	"uint8":     token.Uint8,
	"int16":     token.Int16,
	"uint16":    token.Uint16,
	"struct":    token.Struct,
	"service":   token.Service,
	"enum":      token.Enum,
	"binary":    token.Binary,
	"bool":      token.Bool,
	"string":    token.String,
	"exception": token.Exception,
	"const":     token.Const,
	"map":       token.Map,
}

// Get the token type for an identifier.
func IdentifierTokenType(input string) token.TokenType {
	if t, ok := keywordMap[input]; ok {
		return t
	}

	return token.Identifier
}
