package parser

import (
	"entangle/declarations"
	"entangle/token"
	"fmt"
)

// Parse a type.
func (p *sourceParser) parseType(declarationDesc, self string) (decl declarations.Type, err error) {
	// First, check if we're htting a '*' indicating that the type is nilable.
	nilable := false

	if p.tok.Type == token.TokenType('*') {
		nilable = true

		if err = p.next(); err != nil {
			return
		}
	}

	// Parse the type itself.
	switch p.tok.Type {
	case token.Identifier:
		// The identifier will refer either to an enum or a struct.
		if p.tok.StringValue == self && !nilable {
			err = p.parseErrorHere("non-nilable self references are not allowed")
		} else if structDecl, ok := p.decl.Structs[p.tok.StringValue]; ok {
			decl = declarations.NewStructType(structDecl, nilable)
		} else if enumDecl, ok := p.decl.Enums[p.tok.StringValue]; ok {
			decl = declarations.NewEnumType(enumDecl, nilable)
		} else {
			err = p.parseErrorHere(fmt.Sprintf("unknown type '%s'", p.tok.StringValue))
		}

	case token.Map:
		var keyType, valueType declarations.Type

		// Make sure the next token is a matching '['.
		if err = p.next(); err != nil {
			return
		}

		switch p.tok.Type {
		case token.TokenType('['):
			break

		case token.NewLine:
			err = p.parseErrorHere(fmt.Sprintf("unexpected end of line in %s", declarationDesc))

		case token.EndOfFile:
			err = p.parseErrorHere(fmt.Sprintf("unexpected end of file in %s", declarationDesc))

		default:
			err = p.parseErrorHere("expected '['")
		}

		if err != nil {
			return
		}

		if err = p.next(); err != nil {
			return
		}

		// Parse the key type.
		if keyType, err = p.parseType(declarationDesc, ""); err != nil {
			return
		}

		// Make sure the next token is a matching ']'.
		if err = p.next(); err != nil {
			return
		}

		switch p.tok.Type {
		case token.TokenType(']'):
			break

		case token.NewLine:
			err = p.parseErrorHere(fmt.Sprintf("unexpected end of line in %s", declarationDesc))

		case token.EndOfFile:
			err = p.parseErrorHere(fmt.Sprintf("unexpected end of file in %s", declarationDesc))

		default:
			err = p.parseErrorHere("expected ']'")
		}

		if err != nil {
			return
		}

		if err = p.next(); err != nil {
			return
		}

		// Parse the value type.
		if valueType, err = p.parseType(declarationDesc, ""); err != nil {
			return
		}

		// Return a map type.
		decl = declarations.NewMapType(keyType, valueType, nilable)

	case token.TokenType('['):
		// Make sure the next token is a matching ']'.
		if err = p.next(); err != nil {
			return
		}

		switch p.tok.Type {
		case token.TokenType(']'):
			break

		case token.NewLine:
			err = p.parseErrorHere(fmt.Sprintf("unexpected end of line in %s", declarationDesc))

		case token.EndOfFile:
			err = p.parseErrorHere(fmt.Sprintf("unexpected end of file in %s", declarationDesc))

		default:
			err = p.parseErrorHere("expected ']'")
		}

		if err != nil {
			return
		}

		if err = p.next(); err != nil {
			return
		}

		// Parse the element type.
		var elementType declarations.Type
		if elementType, err = p.parseType(declarationDesc, ""); err != nil {
			return
		}

		// Returns a list type.
		decl = declarations.NewListType(elementType, nilable)

	case token.Bool:
		if nilable {
			return declarations.NilableBoolType, nil
		} else {
			return declarations.BoolType, nil
		}

	case token.String:
		if nilable {
			return declarations.NilableStringType, nil
		} else {
			return declarations.StringType, nil
		}

	case token.Binary:
		if nilable {
			return declarations.NilableBinaryType, nil
		} else {
			return declarations.BinaryType, nil
		}

	case token.Float32:
		if nilable {
			return declarations.NilableFloat32Type, nil
		} else {
			return declarations.Float32Type, nil
		}

	case token.Float64:
		if nilable {
			return declarations.NilableFloat64Type, nil
		} else {
			return declarations.Float64Type, nil
		}

	case token.Int8:
		if nilable {
			return declarations.NilableInt8Type, nil
		} else {
			return declarations.Int8Type, nil
		}

	case token.Int16:
		if nilable {
			return declarations.NilableInt16Type, nil
		} else {
			return declarations.Int16Type, nil
		}

	case token.Int32:
		if nilable {
			return declarations.NilableInt32Type, nil
		} else {
			return declarations.Int32Type, nil
		}

	case token.Int64:
		if nilable {
			return declarations.NilableInt64Type, nil
		} else {
			return declarations.Int64Type, nil
		}

	case token.Uint8:
		if nilable {
			return declarations.NilableUint8Type, nil
		} else {
			return declarations.Uint8Type, nil
		}

	case token.Uint16:
		if nilable {
			return declarations.NilableUint16Type, nil
		} else {
			return declarations.Uint16Type, nil
		}

	case token.Uint32:
		if nilable {
			return declarations.NilableUint32Type, nil
		} else {
			return declarations.Uint32Type, nil
		}

	case token.Uint64:
		if nilable {
			return declarations.NilableUint64Type, nil
		} else {
			return declarations.Uint64Type, nil
		}

	case token.NewLine:
		err = p.parseErrorHere(fmt.Sprintf("unexpected end of line in %s", declarationDesc))

	case token.EndOfFile:
		err = p.parseErrorHere(fmt.Sprintf("unexpected end of file in %s", declarationDesc))

	default:
		err = p.parseErrorHere(fmt.Sprintf("expected type in %s", declarationDesc))
	}

	return
}
