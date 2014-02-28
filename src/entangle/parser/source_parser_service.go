package parser

import (
	"entangle/declarations"
	"entangle/token"
)

// Parse a service declaration.
func (p *sourceParser) parseService() (err error) {
	contextDesc := "service declaration"

	// Parse the name.
	var name string
	var parentDecl *declarations.Service

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
			return p.parseErrorHeref("service name '%s' would override previous type declaration", name)
		}

	default:
		return p.parseErrorHere("expected service name")
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
			return p.parseErrorHeref("unexpected end of file in %s", contextDesc)

		case token.Identifier:
			var found bool
			parentName := p.tok.StringValue

			if parentDecl, found = p.decl.Services[parentName]; !found {
				return p.parseErrorHeref("unknown parent service '%s'", parentName)
			}

		default:
			return p.parseErrorHere("expected parent service name")
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
	var decl *declarations.Service
	documentation := p.documentationParagraphs()

	if parentDecl != nil {
		decl = parentDecl.Inherit(name, documentation)
	} else {
		decl = declarations.NewService(name, documentation)
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

		// Parse the function.
		var function *declarations.Function
		if function, err = p.parseServiceFunction(decl); err != nil {
			return
		}

		decl.AddFunction(function)
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
	p.decl.AddService(decl)

	return p.next()
}

// Parse a service function declaration.
func (p *sourceParser) parseServiceFunction(serviceDecl *declarations.Service) (decl *declarations.Function, err error) {
	contextDesc := "service function definition"
	argumentContextDesc := "service function argument declaration"

	// Parse the name.
	var name string

	switch p.tok.Type {
	case token.NewLine:
		return nil, p.parseErrorHeref("unexpected end of line in %s", contextDesc)

	case token.EndOfFile:
		return nil, p.parseErrorHeref("unexpected end of file in %s", contextDesc)

	case token.Identifier:
		if err = p.validateFunctionName(&p.tok); err != nil {
			return
		}

		name = p.tok.StringValue

		if serviceDecl.FunctionNameInUse(name) {
			return nil, p.parseErrorHeref("function name '%s' has already been declared", name)
		}

	default:
		return nil, p.parseErrorHere("expected function name")
	}

	if err = p.next(); err != nil {
		return
	}

	// Create the function declaration.
	decl = declarations.NewFunction(name, p.documentationParagraphs())

	// The name should be followed by an opening parenthesis ('(').
	if err = p.expectRune('(', contextDesc); err != nil {
		return
	}

	// Parse arguments.
	for {
		if err = p.skipNewLines(); err != nil {
			return
		}

		// If we've reached a '}' here, let's break out of the loop.
		if p.tok.Type == token.TokenType(')') {
			break
		}

		if err = p.skipNewLinesStoreDocumentation(); err != nil {
			return
		}

		// We should have an unsigned integer constant at this point.
		var index uint

		switch p.tok.Type {
		case token.UintConstant:
			index = uint(p.tok.UintValue)

			if index == 0 {
				return nil, p.parseErrorHere("argument indexes are 1-based")
			} else if decl.ArgumentIndexInUse(index) {
				return nil, p.parseErrorHeref("argument index %d already in use", index)
			}

		case token.EndOfFile:
			return nil, p.parseErrorHeref("unexpected end of file in %s", contextDesc)

		default:
			return nil, p.parseErrorHeref("expected argument index in %s", contextDesc)
		}

		if err = p.next(); err != nil {
			return
		}

		// The unsigned integer should be followed by a colon (':').
		if err = p.expectRune(':', argumentContextDesc); err != nil {
			return
		}

		// We should then get a name.
		var name string

		switch p.tok.Type {
		case token.Identifier:
			if err = p.validateArgumentName(&p.tok); err != nil {
				return
			}

			name = p.tok.StringValue

			if decl.ArgumentNameInUse(name) {
				return nil, p.parseErrorHeref("argument named '%s' already declared", name)
			}

		case token.NewLine:
			return nil, p.parseErrorHeref("unexpected end of line in %s", argumentContextDesc)

		case token.EndOfFile:
			return nil, p.parseErrorHeref("unexpected end of file in %s", argumentContextDesc)

		default:
			return nil, p.parseErrorHeref("expected argument name in %s", argumentContextDesc)
		}

		if err = p.next(); err != nil {
			return
		}

		// Then a type.
		var argumentType declarations.Type
		if argumentType, err = p.parseType(argumentContextDesc, decl.Name); err != nil {
			return
		}

		if err = p.next(); err != nil {
			return
		}

		// If this is not the last argument, we should expect a comma here.
		//
		// Mostly we require this before the user throws in a new line, so that
		// we can avoid ever seeing anyone use the horrible Node.js style of
		// "comma on new line" approach. It's gross.
		if p.tok.Type != token.TokenType(')') {
			if err = p.expectRune(',', contextDesc); err != nil {
				return
			}
		}

		// Add the argument to the function declaration.
		decl.AddArgument(index, name, argumentType)
	}

	// The argument list should be followed by a closing parenthesis (')').
	if err = p.expectRune(')', contextDesc); err != nil {
		return
	}

	// Parse the return type.
	if p.tok.Type != token.NewLine && p.tok.Type != token.EndOfFile {
		if decl.ReturnType, err = p.parseType(contextDesc, decl.Name); err != nil {
			return
		}

		if err = p.next(); err != nil {
			return
		}
	}

	// The function declaration should be followed by a new line.
	switch p.tok.Type {
	case token.NewLine:
		break

	case token.EndOfFile:
		return nil, p.parseErrorHeref("unexpected end of file in %s", contextDesc)

	default:
		return nil, p.parseErrorHeref("expected new line after %s", contextDesc)
	}

	err = p.next()
	return
}
