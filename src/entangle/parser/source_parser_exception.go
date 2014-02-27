package parser

import (
	"entangle/declarations"
	"entangle/token"
)

// Parse an exception declaration.
func (p *sourceParser) parseException() (err error) {
	if err = p.next(); err != nil {
		return
	}

	// Parse the name.
	var name string

	switch p.tok.Type {
	case token.NewLine:
		return p.parseErrorHere("unexpected end of line in exception declaration")

	case token.EndOfFile:
		return p.parseErrorHere("unexpected end of file in exception declaration")

	case token.Identifier:
		if err = p.validateTypeName(&p.tok); err != nil {
			return
		}

		name = p.tok.StringValue

		if p.decl.NameInUse(name) {
			return p.parseErrorHeref("exception name '%s' would override previous type declaration", name)
		}

	default:
		return p.parseErrorHere("expected struct name")
	}

	if err = p.next(); err != nil {
		return
	}

	// The name should be followed by an end of line or file.
	switch p.tok.Type {
	case token.NewLine, token.EndOfFile:
		break

	default:
		return p.parseErrorHere("expected new line following exception declaration")
	}

	// Create and add the exception declaration.
	p.decl.AddException(declarations.NewException(name, p.documentationParagraphs()))

	return p.next()
}
