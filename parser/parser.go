// Package parser provides the parser for the Entangle IDL
package parser

import (
	"../declarations"
	"../errors"
	"../lexer"
	"../source"
	"../token"
)

// Parse an Entangle IDL file to an interface declaration.
func Parse(src *source.Source) (interfaceDeclaration *declarations.Interface, err error) {
	interfaceDeclaration = declarations.NewInterface()

	// Create a source parser.
	errorFrames := []errors.ParseErrorFrame{}

	p := &sourceParser{
		lex:                lexer.NewLexer(src, errorFrames),
		src:                src,
		documentationLines: []token.Token{},
		errorFrames:        errorFrames,
		decl:               interfaceDeclaration,
	}

	err = p.parse()
	return
}
