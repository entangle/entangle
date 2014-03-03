package parser

import (
	"entangle/declarations"
	"entangle/errors"
	"entangle/lexer"
	"entangle/source"
	"entangle/token"
)

// Source parser.
//
// Internal context manager for parsing one source file.
type sourceParser struct {
	src                *source.Source
	lex                *lexer.Lexer
	documentationLines []token.Token
	errorFrames        []errors.ParseErrorFrame
	prev               token.Token
	tok                token.Token
	decl               *declarations.Interface
}

func (p *sourceParser) next() (err error) {
	p.prev = p.tok
	p.tok, err = p.lex.Lex()
	return
}

func (p *sourceParser) parse() (err error) {
	// Start out by reading the definition name.
	for {
		if err = p.next(); err != nil {
			return
		}

		switch p.tok.Type {
		case token.NewLine, token.DocumentationLine:
			break

		case token.Definition:
			err = p.parseDefinition()

		case token.EndOfFile:
			err = p.parseErrorHere("unexpected end of file in definition file, expected 'definition'")

		default:
			// Unexpected.
			err = p.parseErrorHere("unexpected token, expected 'definition'")
		}

		if err != nil {
			return
		}

		if p.decl.Name != "" {
			break
		}
	}

	// Read through till the end.
	for p.tok.Type != token.EndOfFile {
		if err = p.next(); err != nil {
			return
		}

		switch p.tok.Type {
		case token.NewLine:
			// Reset the documentation cache if the previous token was not
			// a documentation line.
			if p.prev.Type != token.DocumentationLine {
				p.documentationLines = []token.Token{}
			}

		case token.DocumentationLine:
			// Store the documentation line.
			p.documentationLines = append(p.documentationLines, p.tok)

		case token.Import:
			err = p.parseImport()

		case token.Struct:
			err = p.parseStruct()

		case token.Exception:
			err = p.parseException()

		case token.Enum:
			err = p.parseEnum()

		case token.Service:
			err = p.parseService()

		default:
			// Unexpected.
			p.parseErrorHere("unexpected token")
		}

		if err != nil {
			return
		}
	}

	return
}
