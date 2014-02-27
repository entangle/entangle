package parser

import (
	"entangle/token"
	"strings"
)

func (p *sourceParser) parseImport() (err error) {
	start := p.tok.Start

	// If the next token is an identifier, this is a named import.
	if err = p.next(); err != nil {
		return
	}

	var importName string

	if p.tok.Type == token.Identifier {
		if err = p.validateImportName(&p.tok); err != nil {
			return
		}

		importName = p.tok.StringValue

		if err = p.next(); err != nil {
			return
		}
	}

	// At this point we should have the import path as a literal.
	var path string

	switch p.tok.Type {
	case token.NewLine:
		return p.parseErrorHere("unexpected end of line in import statement")

	case token.EndOfFile:
		return p.parseErrorHere("unexpected end of file in import statement")

	case token.Literal:
		path = strings.TrimSpace(p.tok.StringValue)

		if len(path) == 0 {
			return p.parseErrorHere("empty import path")
		}

	default:
		if len(importName) > 0 {
			return p.parseErrorHere("expected import path")
		} else {
			return p.parseErrorHere("expected import name or import path")
		}
	}

	if err = p.next(); err != nil {
		return
	}

	return p.parseError("imports are currently not supported", start, p.tok.End)
}
