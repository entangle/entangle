package token

import (
	"fmt"
)

// Token type.
type TokenType int64

// List of special tokens.
const (
	/**
	 * String tokens.
	 */
	NewLine TokenType = iota + 0xffffffff
	EndOfFile
	Identifier
	Literal

	/**
	 * Documentation tokens.
	 */
	DocumentationLine

	/**
	 * Header tokens.
	 */
	Import

	/**
	 * Constants.
	 */
	IntConstant
	UintConstant
	FloatConstant

	/**
	 * Base data type tokens.
	 */
	Bool
	String
	Binary
	Float32
	Float64
	Int8
	Int16
	Int32
	Int64
	Uint8
	Uint16
	Uint32
	Uint64

	/**
	 * Definition tokens.
	 */
	Definition
	Const
	Enum
	Struct
	Typedef
	Service
	Exception

	/**
	 * Complex data type tokens.
	 */
	Map
)

var tokenTypeName = map[TokenType]string{
	NewLine:           "NewLine",
	EndOfFile:         "EndOfFile",
	Identifier:        "Identifier",
	Literal:           "Literal",
	DocumentationLine: "DocumentationLine",
	Import:            "Import",
	IntConstant:       "IntConstant",
	UintConstant:      "UintConstant",
	FloatConstant:     "FloatConstant",
	Bool:              "Bool",
	String:            "String",
	Binary:            "Binary",
	Float32:           "Float32",
	Float64:           "Float64",
	Int8:              "Int8",
	Int16:             "Int16",
	Int32:             "Int32",
	Int64:             "Int64",
	Uint8:             "Uint8",
	Uint16:            "Uint16",
	Uint32:            "Uint32",
	Uint64:            "Uint64",
	Definition:        "Definition",
	Const:             "Const",
	Enum:              "Enum",
	Struct:            "Struct",
	Typedef:           "Typedef",
	Service:           "Service",
	Exception:         "Exception",
	Map:               "Map",
}

var tokenTypeRepresentation = map[TokenType]string{
	NewLine:           "new line",
	EndOfFile:         "EOF",
	Identifier:        "identifier",
	Literal:           "literal",
	DocumentationLine: "documentation line",
	Import:            "import",
	IntConstant:       "integer constant",
	UintConstant:      "unsigned integer constant",
	FloatConstant:     "floating point constant",
	Bool:              "bool",
	String:            "string",
	Binary:            "binary",
	Float32:           "float32",
	Float64:           "float64",
	Int8:              "int8",
	Int16:             "int16",
	Int32:             "int32",
	Int64:             "int64",
	Uint8:             "uint8",
	Uint16:            "uint16",
	Uint32:            "uint32",
	Uint64:            "uint64",
	Definition:        "definition",
	Const:             "const",
	Enum:              "enum",
	Struct:            "struct",
	Typedef:           "typedef",
	Service:           "service",
	Exception:         "exception",
	Map:               "map",
}

func (t TokenType) String() string {
	if name, ok := tokenTypeName[t]; ok {
		return name
	}

	if t >= TokenType(33) && t <= TokenType(126) {
		return string(t)
	}

	return fmt.Sprintf("%#U", t)
}

func (t TokenType) Representation() string {
	if name, ok := tokenTypeRepresentation[t]; ok {
		return name
	}

	if t >= TokenType(33) && t <= TokenType(126) {
		return fmt.Sprintf("'%s'", string(t))
	}

	return fmt.Sprintf("%#U", t)
}
