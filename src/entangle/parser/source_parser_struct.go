package parser

import (
	"entangle/declarations"
	"entangle/token"
)

// Parse a struct declaration.
func (p *sourceParser) parseStruct() (err error) {
	contextDesc := "struct declaration"
	fieldContextDesc := "struct field declaration"

	// Parse the name.
	var name string
	var parentDecl *declarations.Struct

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
			return p.parseErrorHeref("struct name '%s' would override previous type declaration", name)
		}

	default:
		return p.parseErrorHere("expected struct name")
	}

	// Skip new lines.
	if err = p.nextAndSkipNewLines(); err != nil {
		return
	}

	// If we get a ':', this is inheriting from another struct.
	if p.tok.Type == token.TokenType(':') {
		if err = p.nextAndSkipNewLines(); err != nil {
			return
		}

		switch p.tok.Type {
		case token.EndOfFile:
			return p.parseErrorHere("unexpected end of file in struct declaration")

		case token.Identifier:
			var found bool
			parentName := p.tok.StringValue

			if parentDecl, found = p.decl.Structs[parentName]; !found {
				return p.parseErrorHeref("unknown parent struct '%s'", parentName)
			}

		default:
			return p.parseErrorHere("expected parent struct name")
		}

		if err = p.nextAndSkipNewLines(); err != nil {
			return
		}
	}

	// At this point, we should be met with an opening curly brace.
	if err = p.expectRune('{', contextDesc); err != nil {
		return
	}

	// Let's initialize the declaration at this point.
	var decl *declarations.Struct
	documentation := p.documentationParagraphs()

	if parentDecl != nil {
		decl = parentDecl.Inherit(name, documentation)
	} else {
		decl = declarations.NewStruct(name, documentation)
	}

	// From here on out, we should be getting documentation and field
	// definitions.
	for {
		if err = p.skipNewLinesStoreDocumentation(); err != nil {
			return
		}

		// If we've reached a '}' here, let's break out of the loop.
		if p.tok.Type == token.TokenType('}') {
			break
		}

		// We should have an unsigned integer constant at this point.
		var fieldIndex uint

		switch p.tok.Type {
		case token.UintConstant:
			fieldIndex = uint(p.tok.UintValue)

			if fieldIndex == 0 {
				return p.parseErrorHere("field indexes are 1-based")
			} else if decl.FieldIndexInUse(fieldIndex) {
				return p.parseErrorHeref("field index %d already in use", fieldIndex)
			}

		case token.EndOfFile:
			return p.parseErrorHere("unexpected end of file in struct declaration")

		default:
			return p.parseErrorHere("expected field index")
		}

		if err = p.next(); err != nil {
			return
		}

		// The unsigned integer should be followed by a colon (':').
		if err = p.expectRune(':', fieldContextDesc); err != nil {
			return
		}

		// We should then get a name.
		var name string

		switch p.tok.Type {
		case token.Identifier:
			if err = p.validateFieldName(&p.tok); err != nil {
				return
			}

			name = p.tok.StringValue

			if decl.FieldNameInUse(name) {
				return p.parseErrorHeref("field name '%s' already in use", name)
			}

		case token.NewLine:
			return p.parseErrorHere("unexpected end of line in struct field declaration")

		case token.EndOfFile:
			return p.parseErrorHere("unexpected end of file in struct field declaration")

		default:
			return p.parseErrorHere("expected field name in struct field definition")
		}

		if err = p.next(); err != nil {
			return
		}

		// Then a type.
		var fieldType declarations.Type
		if fieldType, err = p.parseType(fieldContextDesc, decl.Name); err != nil {
			return
		}

		if err = p.next(); err != nil {
			return
		}

		// And then a new line.
		switch p.tok.Type {
		case token.NewLine:
			break

		case token.EndOfFile:
			return p.parseErrorHere("unexpected end of file in struct declaration")

		default:
			return p.parseErrorHere("expected new line after struct field definition")
		}

		// Add the field to the struct declaration.
		decl.AddField(fieldIndex, name, p.documentationParagraphs(), fieldType)

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
	p.decl.AddStruct(decl)

	return p.next()
}
