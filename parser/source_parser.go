package parser

import (
	"../declarations"
	"../errors"
	"../lexer"
	"../source"
	"../token"
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
