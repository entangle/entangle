package parser

import (
	"entangle/declarations"
	"entangle/token"
	"math"
)

// Parse an enumeration declaration.
func (p *sourceParser) parseEnum() (err error) {
	contextDesc := "enumeration declaration"
	valueContextDesc := "enumeration value declaration"

	// Parse the name.
	var name string

	if err = p.next(); err != nil {
		return
	}

	switch p.tok.Type {
	case token.NewLine:
		return p.parseErrorHeref("unexpected end of line in %s", contextDesc)

	case token.EndOfFile:
		return p.parseErrorHeref("unexpected end of file in %s", contextDesc)

	case token.Identifier:
		if err = p.validateTypeName(&p.tok); err != nil {
			return
		}

		name = p.tok.StringValue

		if p.decl.NameInUse(name) {
			return p.parseErrorHeref("enumeration name '%s' would override previous type declaration", name)
		}

	default:
		return p.parseErrorHere("expected enumeration name")
	}

	// Skip new lines.
	if err = p.nextAndSkipNewLines(); err != nil {
		return
	}

	// At this point, we should be met with an opening curly brace.
	if err = p.expectRune('{', contextDesc); err != nil {
		return
	}

	// Let's initialize the declaration at this point.
	decl := declarations.NewEnum(name, p.documentationParagraphs())
	p.decl.MarkNameAsUsed(name)

	// From here on out, we should be getting documentation and value
	// definitions.
	for {
		if err = p.skipNewLinesStoreDocumentation(); err != nil {
			return
		}

		// If we've reached a '}' here, let's break out of the loop.
		if p.tok.Type == token.TokenType('}') {
			break
		}

		// We should have an integer constant at this point.
		var value int64

		switch p.tok.Type {
		case token.UintConstant:
			if p.tok.UintValue > math.MaxInt64 {
				return p.parseErrorHere("enumeration value out of range")
			}

			value = int64(p.tok.UintValue)

		case token.IntConstant:
			value = p.tok.IntValue

		case token.EndOfFile:
			return p.parseErrorHeref("unexpected end of file in %s", contextDesc)

		default:
			return p.parseErrorHere("expected field index")
		}

		if enumValue, exists := decl.Values[value]; exists {
			return p.parseErrorHeref("another enumeration value in '%s' already has this value: '%s'", name, enumValue.Name)
		}

		if err = p.next(); err != nil {
			return
		}

		// The unsigned integer should be followed by a colon (':').
		if err = p.expectRune(':', valueContextDesc); err != nil {
			return
		}

		// We should then get a name.
		var name string

		switch p.tok.Type {
		case token.Identifier:
			if err = p.validateEnumValueName(&p.tok); err != nil {
				return
			}

			name = p.tok.StringValue

			if p.decl.NameInUse(name) {
				return p.parseErrorHeref("enumeration value name '%s' would override previous type definition", name)
			}

		case token.NewLine:
			return p.parseErrorHere("unexpected end of line in enumeration value declaration")

		case token.EndOfFile:
			return p.parseErrorHere("unexpected end of file in enumeration value declaration")

		default:
			return p.parseErrorHere("expected name in enumeration value declaration")
		}

		if err = p.next(); err != nil {
			return
		}

		// And then a new line.
		switch p.tok.Type {
		case token.NewLine:
			break

		case token.EndOfFile:
			return p.parseErrorHere("unexpected end of file in enumeration declaration")

		default:
			return p.parseErrorHere("expected new line after enumeration value definition")
		}

		// Add the field to the struct declaration.
		decl.AddValue(value, name, p.documentationParagraphs())

		if err = p.next(); err != nil {
			return
		}
	}

	// Here, we should be met with a closing curly brace and a new line or
	// end of file.
	if err = p.expectRune('}', contextDesc); err != nil {
		return
	}

	switch p.tok.Type {
	case token.NewLine, token.EndOfFile:
		break

	default:
		return p.parseErrorHere("expected new line following '}'")
	}

	// Add the declaration to the interface declaration.
	p.decl.AddEnum(decl)

	return p.next()
}
