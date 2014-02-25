package parser

import (
	"../errors"
	"../token"
	"fmt"
)

func (p *sourceParser) parseError(description string, start, end token.Position) error {
	return errors.NewParseError(description, start, end, p.src, p.errorFrames)
}

func (p *sourceParser) parseErrorHere(description string) error {
	return errors.NewParseErrorForToken(description, &p.tok, p.src, p.errorFrames)
}

func (p *sourceParser) parseErrorHeref(description string, a ...interface{}) error {
	return errors.NewParseErrorForToken(fmt.Sprintf(description, a...), &p.tok, p.src, p.errorFrames)
}

func (p *sourceParser) parseErrorForToken(description string, tok *token.Token) error {
	return errors.NewParseErrorForToken(description, tok, p.src, p.errorFrames)
}
