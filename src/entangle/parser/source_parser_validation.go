package parser

import (
	"entangle/token"
	"fmt"
	"regexp"
)

var (
	lowerCamelCaseExpression             = regexp.MustCompile(`^[a-z][a-z0-9]*$`)
	firstLowerCamelCaseExpression        = regexp.MustCompile(`^[a-z][a-zA-Z0-9]*$`)
	firstUpperCamelCaseExpression        = regexp.MustCompile(`^[A-Z][a-zA-Z0-9]*$`)
	firstLowerCamelOrSnakeCaseExpression = regexp.MustCompile(`^[a-z][_a-zA-Z0-9]*$`)
	firstUpperCamelOrSnakeCaseExpression = regexp.MustCompile(`^[A-Z][_a-zA-Z0-9]*$`)
)

// Reserved identifiers.
var reservedIdentifiers = map[string]struct{}{
	"int64":   struct{}{},
	"uint64":  struct{}{},
	"float64": struct{}{},
	"int32":   struct{}{},
	"uint32":  struct{}{},
	"float32": struct{}{},
	"import":  struct{}{},
	"typedef": struct{}{},
	"int8":    struct{}{},
	"uint8":   struct{}{},
	"int16":   struct{}{},
	"uint16":  struct{}{},
	"struct":  struct{}{},
	"service": struct{}{},
	"enum":    struct{}{},
	"binary":  struct{}{},
	"bool":    struct{}{},
	"const":   struct{}{},
}

var reservedArgumentNames = map[string]struct{}{
	"notify": struct{}{},
	"trace":  struct{}{},
}

var reservedFieldNames = map[string]struct{}{
	"Serialize":   struct{}{},
	"Deserialize": struct{}{},
}

var reservedDefinitionNames = map[string]struct{}{
	"entangle": struct{}{},
}

// Validate an import name.
func (p *sourceParser) validateImportName(tok *token.Token) error {
	if _, reserved := reservedIdentifiers[tok.StringValue]; reserved {
		return p.parseErrorForToken(fmt.Sprintf("'%s' is a reserved identifier", tok.StringValue), tok)
	}

	if !firstLowerCamelOrSnakeCaseExpression.MatchString(tok.StringValue) {
		return p.parseErrorForToken(fmt.Sprintf("'%s' is not a valid import name. Import names must be lower camel case or lower snake case", tok.StringValue), tok)
	}

	return nil
}

// Validate a type name.
func (p *sourceParser) validateTypeName(tok *token.Token) error {
	if _, reserved := reservedIdentifiers[tok.StringValue]; reserved {
		return p.parseErrorForToken(fmt.Sprintf("'%s' is a reserved identifier", tok.StringValue), tok)
	}

	if !firstUpperCamelCaseExpression.MatchString(tok.StringValue) {
		return p.parseErrorForToken(fmt.Sprintf("'%s' is not a valid type name. Type names must be upper camel case", tok.StringValue), tok)
	}

	return nil
}

// Validate a function name.
func (p *sourceParser) validateFunctionName(tok *token.Token) error {
	if _, reserved := reservedIdentifiers[tok.StringValue]; reserved {
		return p.parseErrorForToken(fmt.Sprintf("'%s' is a reserved identifier", tok.StringValue), tok)
	}

	if !firstUpperCamelCaseExpression.MatchString(tok.StringValue) {
		return p.parseErrorForToken(fmt.Sprintf("'%s' is not a valid function name. Function names must be upper camel case", tok.StringValue), tok)
	}

	return nil
}

// Validate an enumeration value name.
func (p *sourceParser) validateEnumValueName(tok *token.Token) error {
	if _, reserved := reservedIdentifiers[tok.StringValue]; reserved {
		return p.parseErrorForToken(fmt.Sprintf("'%s' is a reserved identifier", tok.StringValue), tok)
	}

	if !firstUpperCamelCaseExpression.MatchString(tok.StringValue) {
		return p.parseErrorForToken(fmt.Sprintf("'%s' is not a valid enumeration value name. Enumeration value names must be upper camel case or upper snake case", tok.StringValue), tok)
	}

	return nil
}

// Validate a field name.
func (p *sourceParser) validateFieldName(tok *token.Token) error {
	if _, reserved := reservedIdentifiers[tok.StringValue]; reserved {
		return p.parseErrorForToken(fmt.Sprintf("'%s' is a reserved identifier", tok.StringValue), tok)
	}

	if _, reserved := reservedFieldNames[tok.StringValue]; reserved {
		return p.parseErrorForToken(fmt.Sprintf("'%s' is a reserved field name", tok.StringValue), tok)
	}

	if !firstUpperCamelCaseExpression.MatchString(tok.StringValue) {
		return p.parseErrorForToken(fmt.Sprintf("'%s' is not a valid field name. Field names must be upper camel case", tok.StringValue), tok)
	}

	return nil
}

// Validate an argument name.
func (p *sourceParser) validateArgumentName(tok *token.Token) error {
	if _, reserved := reservedIdentifiers[tok.StringValue]; reserved {
		return p.parseErrorForToken(fmt.Sprintf("'%s' is a reserved identifier", tok.StringValue), tok)
	}

	if _, reserved := reservedArgumentNames[tok.StringValue]; reserved {
		return p.parseErrorForToken(fmt.Sprintf("'%s' is a reserved argument name", tok.StringValue), tok)
	}

	if !firstLowerCamelCaseExpression.MatchString(tok.StringValue) {
		return p.parseErrorForToken(fmt.Sprintf("'%s' is not a valid argument name. Argument names must be lower camel case", tok.StringValue), tok)
	}

	return nil
}

// Validate a definition name.
func (p *sourceParser) validateDefinitionName(tok *token.Token) error {
	if _, reserved := reservedIdentifiers[tok.StringValue]; reserved {
		return p.parseErrorForToken(fmt.Sprintf("'%s' is a reserved identifier", tok.StringValue), tok)
	}

	if _, reserved := reservedDefinitionNames[tok.StringValue]; reserved {
		return p.parseErrorForToken(fmt.Sprintf("'%s' is a reserved definition name", tok.StringValue), tok)
	}

	if !lowerCamelCaseExpression.MatchString(tok.StringValue) {
		return p.parseErrorForToken(fmt.Sprintf("'%s' is not a valid definition name. Definition names must be lower snake case", tok.StringValue), tok)
	}

	return nil
}
