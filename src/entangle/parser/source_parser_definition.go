package parser

import (
	"entangle/token"
)

func (p *sourceParser) parseDefinition() (err error) {
	if err = p.next(); err != nil {
		return
	}

	// We should have the definition name as an identifier.
	switch p.tok.Type {
	case token.NewLine:
		return p.parseErrorHere("unexpected end of line in definition statement")

	case token.EndOfFile:
		return p.parseErrorHere("unexpected end of file in definition statement")

	case token.Identifier:
		if err = p.validateDefinitionName(&p.tok); err != nil {
			return
		}

		p.decl.Name = p.tok.StringValue

	default:
		return p.parseErrorHere("expected definition name")
	}

	if err = p.next(); err != nil {
		return
	}

	// The definition name should be followed by a new line.
	switch p.tok.Type {
	case token.NewLine, token.EndOfFile:
		return p.next()

	default:
		return p.parseErrorHere("expected new line following definition name")
	}
}
