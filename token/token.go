// Package token provides the token interchanged between the lexer and parser.
package token

// Token.
type Token struct {
	// Token type.
	Type TokenType

	// Start position.
	Start Position

	// End position.
	End Position

	// String value.
	StringValue string

	// Signed integer value.
	IntValue int64

	// Unsigned integer value.
	UintValue uint64

	// Floating point value.
	FloatValue float64
}
